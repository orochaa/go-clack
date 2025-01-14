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
