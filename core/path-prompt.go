package core

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/Mist3rBru/go-clack/core/internals"
	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type PathPrompt struct {
	Prompt[string]
	OnlyShowDir bool
	Required    bool
	Hint        string
	HintOptions []string
	HintIndex   int
	FileSystem  FileSystem
}

type PathPromptParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue string
	OnlyShowDir  bool
	Required     bool
	FileSystem   FileSystem
	Validate     func(value string) error
	Render       func(p *PathPrompt) string
}

// NewPathPrompt initializes and returns a new instance of PathPrompt.
//
// The user can input a path.
// The prompt has built-in autosuggestion and autocomplete features.
// The prompt returns the path.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue (string): The initial value of the path input (default: current working directory).
//   - OnlyShowDir (bool): Whether to only show directories (default: false).
//   - Required (bool): Whether the path input is required (default: false).
//   - FileSystem (FileSystem): The file system implementation to use (default: OSFileSystem).
//   - Validate (func(value string) error): Custom validation function for the path (default: nil).
//   - Render (func(p *PathPrompt) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *PathPrompt: A new instance of PathPrompt.
func NewPathPrompt(params PathPromptParams) *PathPrompt {
	v := validator.NewValidator("PathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p PathPrompt
	p = PathPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  len(params.InitialValue),
			Validate:     WrapValidate(params.Validate, &p.Required, "Path does not exist! Please enter a valid path."),
			Render:       WrapRender[string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
		HintIndex:   -1,
		Required:    params.Required,
		FileSystem:  params.FileSystem,
	}

	if cwd, err := p.FileSystem.Getwd(); err == nil && params.InitialValue == "" {
		p.Prompt.Value = cwd
		p.Value = cwd
		p.CursorIndex = len(cwd)
	}
	p.changeHint()

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

// mapHintOptions generates a list of hint options based on the current path value.
// It filters entries in the current directory that match the end of the path value.
//
// Returns:
//   - []string: A slice of hint options.
func (p *PathPrompt) mapHintOptions() []string {
	options := []string{}
	dirPathRegex := regexp.MustCompile(`^(.*)/.*\s*`)
	dirPath := dirPathRegex.ReplaceAllString(p.Value, "$1")

	if strings.HasPrefix(dirPath, "~") {
		if homeDir, err := p.FileSystem.UserHomeDir(); err == nil {
			dirPath = strings.Replace(dirPath, "~", homeDir, 1)
		}
	}

	entries, err := p.FileSystem.ReadDir(dirPath)
	if err != nil {
		return options
	}

	for _, entry := range entries {
		if (p.OnlyShowDir && !entry.IsDir()) || !strings.HasPrefix(entry.Name(), p.valueEnd()) {
			continue
		}

		option := entry.Name()
		if entry.IsDir() {
			option += "/"
		}

		options = append(options, option)
	}

	return options
}

// valueEnd extracts the last segment of the current path value.
// This is used to match hint options with the current input.
//
// Returns:
//   - string: The last segment of the path value.
func (p *PathPrompt) valueEnd() string {
	valueEndRegex := regexp.MustCompile("^.*/(.*)$")
	valueEnd := valueEndRegex.ReplaceAllString(p.Value, "$1")
	return valueEnd
}

// changeHint updates the hint based on the current path value and available hint options.
// If no hint options are available, the hint is cleared.
func (p *PathPrompt) changeHint() {
	hintOptions := p.mapHintOptions()
	p.HintOptions = []string{}
	if len(hintOptions) > 0 {
		p.Hint = strings.Replace(hintOptions[0], p.valueEnd(), "", 1)
	} else {
		p.Hint = ""
	}
}

// ValueWithCursor returns the current path value with a cursor indicator.
// The cursor is represented by an inverse character at the current cursor position.
// If the cursor is at the end of the value, the hint is displayed.
//
// Returns:
//   - string: The path value with the cursor indicator and hint.
func (p *PathPrompt) ValueWithCursor() string {
	var (
		value string
		hint  string
	)
	if p.CursorIndex >= len(p.Value) {
		value = p.Value
		if p.Hint == "" {
			hint = "█"
		} else {
			hint = picocolors.Inverse(string(p.Hint[0])) + picocolors.Dim(p.Hint[1:])
		}
	} else {
		s1 := p.Value[0:p.CursorIndex]
		s2 := p.Value[p.CursorIndex:]
		value = s1 + picocolors.Inverse(string(s2[0])) + s2[1:]
		hint = picocolors.Dim(p.Hint)
	}
	return value + hint
}

// completeValue appends the current hint to the path value and updates the cursor position.
// The hint is then cleared, and new hint options are generated.
func (p *PathPrompt) completeValue() {
	p.Value += p.Hint
	p.Prompt.Value = p.Value
	p.CursorIndex = len(p.Value)
	p.Hint = ""
	p.HintOptions = []string{}
	p.changeHint()
}

// tabComplete handles tab completion for the path input.
// If there is only one hint option, it completes the value.
// Otherwise, it cycles through the available hint options.
func (p *PathPrompt) tabComplete() {
	hintOption := p.mapHintOptions()
	if len(hintOption) == 1 {
		p.completeValue()
	} else if len(p.HintOptions) == 0 {
		p.HintOptions = hintOption
		p.HintIndex = 0
	} else {
		p.HintIndex = utils.MinMaxIndex(p.HintIndex+1, len(p.HintOptions))
		p.Hint = strings.Replace(p.HintOptions[p.HintIndex], p.valueEnd(), "", 1)
	}
}

// handleKeyPress processes key events for the path input.
// It updates the path value and cursor position based on the key pressed.
// Special keys like Tab and Right Arrow trigger hint completion.
//
// Parameters:
//   - key (*Key): The key event to process.
func (p *PathPrompt) handleKeyPress(key *Key) {
	p.Value, p.CursorIndex = p.TrackKeyValue(key, p.Value, p.CursorIndex)
	if key.Name == RightKey && p.CursorIndex >= len(p.Value) {
		p.completeValue()
	} else if key.Name == TabKey {
		p.tabComplete()
	} else {
		p.changeHint()
	}
}
