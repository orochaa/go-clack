package core

import (
	"context"
	"os"
	"strings"

	"github.com/orochaa/go-clack/core/validator"
	"github.com/orochaa/go-clack/third_party/picocolors"
)

type PasswordPrompt struct {
	Prompt[string]
	Required bool
}

type PasswordPromptParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue string
	Required     bool
	Validate     func(value string) error
	Render       func(p *PasswordPrompt) string
}

// NewPasswordPrompt initializes and returns a new instance of PasswordPrompt.
//
// The user can input a password.
// The password is masked by asterisks ("*").
// The prompt returns the password without the mask.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue (string): The initial value of the password input (default: "").
//   - Required (bool): Whether the password input is required (default: false).
//   - Validate (func(value string) error): Custom validation function for the password (default: nil).
//   - Render (func(p *PasswordPrompt) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *PasswordPrompt: A new instance of PasswordPrompt.
func NewPasswordPrompt(params PasswordPromptParams) *PasswordPrompt {
	v := validator.NewValidator("PasswordPrompt")
	v.ValidateRender(params.Render)

	var p PasswordPrompt
	p = PasswordPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  len(params.InitialValue),
			Validate:     WrapValidate(params.Validate, &p.Required, "Password is required! Please enter a value."),
			Render:       WrapRender[string](&p, params.Render),
		}),
		Required: params.Required,
	}

	p.On(KeyEvent, func(args ...any) {
		p.Value, p.CursorIndex = p.TrackKeyValue(args[0].(*Key), p.Value, p.CursorIndex)
	})

	return &p
}

// ValueWithMask returns the current password value masked with asterisks (*).
// This is useful for displaying the password in a secure manner.
//
// Returns:
//   - string: The masked password value.
func (p *PasswordPrompt) ValueWithMask() string {
	return strings.Repeat("*", len(p.Value))
}

// ValueWithMaskAndCursor returns the current password value masked with asterisks (*) and includes a cursor indicator.
// The cursor is represented by an inverse character at the current cursor position.
//
// Returns:
//   - string: The masked password value with the cursor indicator.
func (p *PasswordPrompt) ValueWithMaskAndCursor() string {
	maskedValue := strings.Repeat("*", len(p.Value))
	if p.CursorIndex == len(p.Value) {
		return maskedValue + "█"
	}
	return maskedValue[0:p.CursorIndex] + picocolors.Inverse(string(maskedValue[p.CursorIndex])) + maskedValue[p.CursorIndex+1:]
}
