package core

import (
	"context"
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
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

func (p *PasswordPrompt) ValueWithMask() string {
	return strings.Repeat("*", len(p.Value))
}

func (p *PasswordPrompt) ValueWithMaskAndCursor() string {
	maskedValue := strings.Repeat("*", len(p.Value))
	if p.CursorIndex == len(p.Value) {
		return maskedValue + "█"
	}
	return maskedValue[0:p.CursorIndex] + picocolors.Inverse(string(maskedValue[p.CursorIndex])) + maskedValue[p.CursorIndex+1:]
}
