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

func (p *TextPrompt) handleKeyPress(key *Key) {
	if key.Name == TabKey && p.Value == "" && p.Placeholder != "" {
		p.Value = p.Placeholder
		p.CursorIndex = len(p.Placeholder)
		return
	}

	p.Value, p.CursorIndex = p.TrackKeyValue(key, p.Value, p.CursorIndex)
}

func (p *TextPrompt) ValueWithCursor() string {
	if p.CursorIndex == len(p.Value) {
		return p.Value + "█"
	}
	return p.Value[0:p.CursorIndex] + picocolors.Inverse(string(p.Value[p.CursorIndex])) + p.Value[p.CursorIndex+1:]
}
