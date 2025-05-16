package core

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/core/validator"
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
}

func (o *SelectOption[TValue]) String() string {
	if o == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"{Label:%s, Value:%v(%T)}",
		o.Label,
		o.Value, o.Value,
	)
}

type SelectPrompt[TValue comparable] struct {
	Prompt[TValue]
	initialOptions []*SelectOption[TValue]
	Options        []*SelectOption[TValue]
	Search         string
	Filter         bool
	Required       bool
}

type SelectPromptParams[TValue comparable] struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue TValue
	Options      []*SelectOption[TValue]
	Filter       bool
	Required     bool
	Render       func(p *SelectPrompt[TValue]) string
}

// NewSelectPrompt initializes and returns a new instance of SelectPrompt.
//
// The user can navigate through options using arrow keys.
// The user can select an option using enter key.
// The prompt returns the value of the selected option.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue (TValue): The initial value of the prompt (default: zero value of TValue).
//   - Options ([]*SelectOption[TValue]): A list of options for the prompt (default: nil).
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - Required (bool): Whether the prompt requires a selection (default: false).
//   - Render (func(p *SelectPrompt[TValue]) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *SelectPrompt[TValue]: A new instance of SelectPrompt.
func NewSelectPrompt[TValue comparable](params SelectPromptParams[TValue]) *SelectPrompt[TValue] {
	v := validator.NewValidator("SelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	startIndex := 0
	for i, option := range params.Options {
		if value, ok := any(option.Value).(string); ok && value == "" {
			option.Value = any(option.Label).(TValue)
		}
		if option.Value == params.InitialValue {
			startIndex = i
		}
	}

	var p SelectPrompt[TValue]
	p = SelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[TValue]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.Options[startIndex].Value,
			CursorIndex:  startIndex,
			Validate:     WrapValidate[TValue](nil, &p.Required, "Please select an option."),
			Render:       WrapRender[TValue](&p, params.Render),
		}),
		initialOptions: params.Options,
		Options:        params.Options,
		Filter:         params.Filter,
		Required:       params.Required,
	}

	actionHandler := NewActionHandler(map[Action]func(){
		UpAction:    func() { p.moveCursor(-1) },
		DownAction:  func() { p.moveCursor(1) },
		LeftAction:  func() { p.moveCursor(-1) },
		RightAction: func() { p.moveCursor(1) },
		HomeAction:  func() { p.CursorIndex = 0 },
		EndAction:   func() { p.CursorIndex = len(p.Options) - 1 },
	}, p.filterOptions)
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))

		if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
			p.Value = p.Options[p.CursorIndex].Value
		} else {
			p.Value = *new(TValue)
		}
	})

	return &p
}

// moveCursor moves the cursor up or down within the list of options.
// It ensures the cursor stays within the bounds of the available options.
//
// Parameters:
//   - direction (int): The direction to move the cursor (-1 for up, 1 for down).
func (p *SelectPrompt[TValue]) moveCursor(direction int) {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
}

// filterOptions updates the search term based on the provided key input and filters the available options.
// If the search term is empty, it resets the options to the initial list.
//
// Parameters:
//   - key (*Key): The key event that triggered the filtering.
func (p *SelectPrompt[TValue]) filterOptions(key *Key) {
	if !p.Filter {
		return
	}

	p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
	p.CursorIndex = 0

	if p.Search == "" {
		p.Options = p.initialOptions
		for i, option := range p.Options {
			if option.Value == p.Value {
				p.CursorIndex = i
				break
			}
		}
		return
	}

	p.Options = []*SelectOption[TValue]{}
	searchRegex, err := regexp.Compile("(?i)" + p.Search)
	if err != nil {
		return
	}

	for _, option := range p.initialOptions {
		if matched := searchRegex.MatchString(option.Label); matched {
			p.Options = append(p.Options, option)
			if option.Value == p.Value {
				p.CursorIndex = len(p.Options) - 1
			}
		}
	}
}
