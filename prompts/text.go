package prompts

import (
	"context"
	"os"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
)

type TextParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	Message      string
	InitialValue string
	Placeholder  string
	Required     bool
	Validate     func(value string) error
}

// Text displays a input prompt to the user.
//
// The prompt displays a message.
// The user can input a value.
// The prompt returns the value.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Message (string): The message to display to the user (default: "").
//   - InitialValue (string): The initial value of the text input (default: "").
//   - Placeholder (string): The placeholder text to display when the input is empty (default: "").
//   - Required (bool): Whether the text input is required (default: false).
//   - Validate (func(value string) error): Custom validation function for the input (default: nil).
//
// Returns:
//   - string: The typed value.
//   - error: An error if the user cancels the prompt or if an error occurs.
func Text(params TextParams) (string, error) {
	p := core.NewTextPrompt(core.TextPromptParams{
		Context:      params.Context,
		Input:        params.Input,
		Output:       params.Output,
		InitialValue: params.InitialValue,
		Placeholder:  params.Placeholder,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.TextPrompt) string {
			return theme.ApplyTheme(theme.ThemeParams[string]{
				Context:         p.Prompt,
				Message:         params.Message,
				Value:           p.Value,
				ValueWithCursor: p.ValueWithCursor(),
				Placeholder:     p.Placeholder,
			})
		},
	})
	test.TextTestingPrompt = p
	return p.Run()
}
