package core

import (
	"context"
	"os"

	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type TextPrompt struct {
	Prompt[string]
	Placeholder string
	Required    bool
}

type TextPromptParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue string
	Placeholder  string
	Required     bool
	Validate     func(value string) error
	Render       func(p *TextPrompt) string
}

// NewTextPrompt initializes and returns a new instance of TextPrompt.
//
// The user can input a value.
// The prompt returns the value.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue (string): The initial value of the text input (default: "").
//   - Placeholder (string): The placeholder text to display when the input is empty (default: "").
//   - Required (bool): Whether the text input is required (default: false).
//   - Validate (func(value string) error): Custom validation function for the input (default: nil).
//   - Render (func(p *TextPrompt) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *TextPrompt: A new instance of TextPrompt.
func NewTextPrompt(params TextPromptParams) *TextPrompt {
	v := validator.NewValidator("TextPrompt")
	v.ValidateRender(params.Render)

	var p TextPrompt
	p = TextPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  len(params.InitialValue),
			Validate:     WrapValidate(params.Validate, &p.Required, "Value is required! Please enter a value."),
			Render:       WrapRender[string](&p, params.Render),
		}),
		Placeholder: params.Placeholder,
		Required:    params.Required,
	}

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

// handleKeyPress processes key press events for the TextPrompt.
// It updates the input value and cursor position based on the key pressed.
// If the Tab key is pressed and the input is empty, the placeholder is inserted as the value.
//
// Parameters:
//   - key (*Key): The key event to process.
func (p *TextPrompt) handleKeyPress(key *Key) {
	if key.Name == TabKey && p.Value == "" && p.Placeholder != "" {
		p.Value = p.Placeholder
		p.CursorIndex = len(p.Placeholder)
		return
	}

	p.Value, p.CursorIndex = p.TrackKeyValue(key, p.Value, p.CursorIndex)
}

// ValueWithCursor returns the current input value with a cursor indicator.
// The cursor is represented by an inverse character at the current cursor position.
// If the cursor is at the end of the value, it is displayed as a block character.
//
// Returns:
//   - string: The input value with the cursor indicator.
func (p *TextPrompt) ValueWithCursor() string {
	if p.CursorIndex == len(p.Value) {
		return p.Value + "█"
	}
	return p.Value[0:p.CursorIndex] + picocolors.Inverse(string(p.Value[p.CursorIndex])) + p.Value[p.CursorIndex+1:]
}
