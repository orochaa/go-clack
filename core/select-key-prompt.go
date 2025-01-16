package core

import (
	"context"
	"os"

	"github.com/Mist3rBru/go-clack/core/validator"
)

type SelectKeyOption[TValue any] struct {
	Label string
	Value TValue
	Key   string
}

type SelectKeyPrompt[TValue any] struct {
	Prompt[TValue]
	Options []*SelectKeyOption[TValue]
}

type SelectKeyPromptParams[TValue any] struct {
	Context context.Context
	Input   *os.File
	Output  *os.File
	Options []*SelectKeyOption[TValue]
	Render  func(p *SelectKeyPrompt[TValue]) string
}

// NewSelectKeyPrompt initializes and returns a new instance of SelectKeyPrompt.
// It sets up the prompt, validates the render function, and configures the options for key-based selection.

// The user can select an option by clicking on the related key.
// The prompt returns the value of the selected option.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Options ([]*SelectKeyOption[TValue]): A list of options for the prompt (default: nil).
//   - Render (func(p *SelectKeyPrompt[TValue]) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *SelectKeyPrompt[TValue]: A new instance of SelectKeyPrompt.
func NewSelectKeyPrompt[TValue any](params SelectKeyPromptParams[TValue]) *SelectKeyPrompt[TValue] {
	v := validator.NewValidator("SelectKeyPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	for _, option := range params.Options {
		if value, ok := any(option.Value).(string); ok && value == "" {
			option.Value = any(option.Key).(TValue)
		}
	}

	var p SelectKeyPrompt[TValue]
	p = SelectKeyPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[TValue]{
			Context: params.Context,
			Input:   params.Input,
			Output:  params.Output,
			Render:  WrapRender[TValue](&p, params.Render),
		}),
		Options: params.Options,
	}

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

// handleKeyPress processes key press events for the SelectKeyPrompt.
// It checks if the pressed key matches any of the option keys and updates the prompt state accordingly.
//
// Parameters:
//   - key (*Key): The key event to process.
func (p *SelectKeyPrompt[TValue]) handleKeyPress(key *Key) {
	for i, option := range p.Options {
		if key.Name == KeyName(option.Key) {
			p.State = SubmitState
			p.Value = option.Value
			p.CursorIndex = i
			return
		}
	}
	if key.Name == EnterKey {
		key.Name = ""
	}
}
