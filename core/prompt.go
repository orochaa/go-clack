package core

import (
	"bufio"
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/core/validator"
	"github.com/orochaa/go-clack/third_party/sisteransi"

	"golang.org/x/term"
)

type State int

const (
	// InitialState is the initial state of the prompt
	InitialState State = iota
	// ActiveState is set after the user's first action
	ActiveState
	// ValidateState is set after 400ms of validation (e.g., checking user input)
	ValidateState
	// ErrorState is set if there is an error during validation
	ErrorState
	// CancelState is set after the user cancels the prompt
	CancelState
	// SubmitState is set after the user submits the input
	SubmitState
)

type Prompt[TValue any] struct {
	context   context.Context
	listeners map[Event][]EventListener

	rl     *bufio.Reader
	input  *os.File
	output *os.File

	State       State
	Error       string
	Value       TValue
	CursorIndex int

	Validate           func(value TValue) error
	ValidationDuration time.Duration
	IsValidating       bool

	Render func(p *Prompt[TValue]) string
	Frame  string
}

type PromptParams[TValue any] struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue TValue
	CursorIndex  int
	Validate     func(value TValue) error
	Render       func(p *Prompt[TValue]) string
}

// NewPrompt initializes a new Prompt with the provided parameters.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue (TValue): The initial value of the prompt (default: zero value of TValue).
//   - CursorIndex (int): The initial cursor position in the input (default: 0).
//   - Validate (func(value TValue) error): Custom validation function for the input (default: nil).
//   - Render (func(p *Prompt[TValue]) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *Prompt[TValue]: A new instance of Prompt.
func NewPrompt[TValue any](params PromptParams[TValue]) *Prompt[TValue] {
	v := validator.NewValidator("Prompt")
	v.ValidateRender(params.Render)

	if params.Context == nil {
		params.Context = context.Background()
	}
	if params.Input == nil {
		params.Input = os.Stdin
	}
	if params.Output == nil {
		params.Output = os.Stdout
	}

	return &Prompt[TValue]{
		context:   params.Context,
		listeners: make(map[Event][]EventListener),

		input:  params.Input,
		output: params.Output,
		rl:     bufio.NewReader(params.Input),

		State:       InitialState,
		Value:       params.InitialValue,
		CursorIndex: params.CursorIndex,

		Validate: params.Validate,
		Render:   params.Render,
	}
}

type KeyName string

type Key struct {
	Name  KeyName
	Char  string
	Shift bool
	Ctrl  bool
}

const (
	UpKey        KeyName = "Up"
	DownKey      KeyName = "Down"
	LeftKey      KeyName = "Left"
	RightKey     KeyName = "Right"
	HomeKey      KeyName = "Home"
	EndKey       KeyName = "End"
	SpaceKey     KeyName = "Space"
	EnterKey     KeyName = "Enter"
	CancelKey    KeyName = "Cancel"
	TabKey       KeyName = "Tab"
	BackspaceKey KeyName = "Backspace"
	EscapeKey    KeyName = "Escape"
)

// ParseKey parses a rune into a Key.
func (p *Prompt[TValue]) ParseKey(r rune) *Key {
	// TODO: parse Backtab(shift+tab) and other variations of shift and ctrl
	switch r {
	case '\r', '\n':
		return &Key{Name: EnterKey}
	case ' ':
		return &Key{Name: SpaceKey}
	case '\b', 127:
		return &Key{Name: BackspaceKey}
	case '\t':
		return &Key{Name: TabKey}
	case 3:
		return &Key{Name: CancelKey}
	case 27:
		readerReady := make(chan bool, 1)
		go func() {
			_, err := p.rl.Peek(2)
			readerReady <- err == nil
		}()

		select {
		case ready := <-readerReady:
			if ready {
				next, err := p.rl.Peek(2)
				if err == nil && len(next) == 2 && next[0] == '[' {
					p.rl.ReadByte() // Consume '['
					thirdByte, _ := p.rl.ReadByte()

					switch thirdByte {
					case 'A':
						return &Key{Name: UpKey}
					case 'B':
						return &Key{Name: DownKey}
					case 'C':
						return &Key{Name: RightKey}
					case 'D':
						return &Key{Name: LeftKey}
					case 'H':
						return &Key{Name: HomeKey}
					case 'F':
						return &Key{Name: EndKey}
					}
				}
				return &Key{}
			} else {
				return &Key{Name: EscapeKey}
			}

		case <-time.After(50 * time.Millisecond):
			return &Key{Name: EscapeKey}
		}
	}

	char := string(r)
	return &Key{Char: char, Name: KeyName(char)}
}

// PressKey handles key press events and updates the state of the prompt.
func (p *Prompt[TValue]) PressKey(key *Key) {
	if p.State == InitialState || p.State == ErrorState {
		p.State = ActiveState
	}

	p.Emit(KeyEvent, key)

	if action, actionExists := Settings.Aliases[key.Name]; actionExists {
		if action == SubmitAction {
			if err := p.validate(); err != nil {
				p.State = ErrorState
				p.Error = err.Error()
			} else {
				p.State = SubmitState
			}
		} else if action == CancelAction {
			p.State = CancelState
		}
	}

	if p.State == SubmitState || p.State == CancelState {
		p.Emit(FinalizeEvent)
	}

	p.render()

	if p.State == SubmitState {
		p.Emit(SubmitEvent)
	} else if p.State == CancelState {
		p.Emit(CancelEvent)
	}
}

// validate performs validation on the current value of the prompt.
func (p *Prompt[TValue]) validate() error {
	if p.Validate == nil {
		return nil
	}

	p.State = ValidateState
	p.IsValidating = true
	p.Emit(ValidateEvent)

	go func() {
		validationStart := time.Now()
		time.Sleep(400 * time.Millisecond)
		for p.IsValidating {
			p.ValidationDuration = time.Since(validationStart)
			p.render()
			time.Sleep(125 * time.Millisecond)
		}
	}()

	err := p.Validate(p.Value)
	p.IsValidating = false

	return err
}

// DiffLines calculates the difference between an old and a new frame.
func (p *Prompt[TValue]) DiffLines(oldFrame, newFrame string) []int {
	var diff []int

	if oldFrame == newFrame {
		return diff
	}

	oldLines := utils.SplitLines(oldFrame)
	newLines := utils.SplitLines(newFrame)
	for i := range max(len(oldLines), len(newLines)) {
		if i >= len(oldLines) || i >= len(newLines) || oldLines[i] != newLines[i] {
			diff = append(diff, i)
		}
	}

	return diff
}

// Size retrieves the width and height of the terminal output.
func (p *Prompt[TValue]) Size() (width int, height int, err error) {
	return term.GetSize(int(p.output.Fd()))
}

// render renders a new frame to the output.
func (p *Prompt[TValue]) render() {
	frame := p.Render(p)

	if p.State == InitialState {
		p.output.WriteString(sisteransi.HideCursor())
		p.output.WriteString(frame)
		p.Frame = frame
		return
	}

	if frame == p.Frame {
		return
	}

	diff := p.DiffLines(frame, p.Frame)
	diffLineIndex := diff[0]
	prevFrameLines := utils.SplitLines((p.Frame))

	// Move to first diff line
	p.output.WriteString(sisteransi.MoveCursor(-(len(prevFrameLines) - 1), -999))
	p.output.WriteString(sisteransi.MoveCursor(diffLineIndex, 0))
	p.output.WriteString(sisteransi.EraseDown())
	lines := utils.SplitLines(frame)
	newLines := lines[diffLineIndex:]
	p.output.WriteString(strings.Join(newLines, "\r\n"))
	p.Frame = frame
}

// Run runs the prompt and processes input.
func (p *Prompt[TValue]) Run() (TValue, error) {
	var oldState *term.State
	if flag.Lookup("test.v") == nil {
		var err error
		oldState, err = term.MakeRaw(int(p.input.Fd()))
		if err != nil {
			return p.Value, err
		}
		defer term.Restore(int(p.input.Fd()), oldState)
	}

	done := make(chan struct{})
	closeCb := func(args ...any) {
		p.output.WriteString(sisteransi.ShowCursor())
		p.output.WriteString("\r\n")
		close(done)
	}
	p.Once(SubmitEvent, closeCb)
	p.Once(CancelEvent, closeCb)

	p.render()

	go func() {
		select {
		case <-done:
			return
		case <-p.context.Done():
			// Restore terminal immediately when context is cancelled
			if oldState != nil {
				term.Restore(int(p.input.Fd()), oldState)
			}
			p.PressKey(&Key{Name: CancelKey})
		}
	}()

outer:
	for {
		select {
		case <-done:
			break outer
		default:
			r, size, err := p.rl.ReadRune()
			if err != nil || size == 0 || p.IsValidating {
				continue
			}
			// Check if cancelled while we were blocked on ReadRune
			select {
			case <-done:
				break outer
			default:
			}
			key := p.ParseKey(r)
			p.PressKey(key)
		}
	}

	if p.State == CancelState {
		return p.Value, ErrCancelPrompt
	}

	return p.Value, nil
}
