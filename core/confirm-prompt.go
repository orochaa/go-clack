package core

import (
	"context"
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type ConfirmPrompt struct {
	Prompt[bool]
	Active   string
	Inactive string
}

type ConfirmPromptParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	Active       string
	Inactive     string
	InitialValue bool
	Render       func(p *ConfirmPrompt) string
}

func NewConfirmPrompt(params ConfirmPromptParams) *ConfirmPrompt {
	v := validator.NewValidator("ConfirmPrompt")
	v.ValidateRender(params.Render)

	if params.Active == "" {
		params.Active = "yes"
	}
	if params.Inactive == "" {
		params.Inactive = "no"
	}

	var p ConfirmPrompt
	p = ConfirmPrompt{
		Prompt: *NewPrompt(PromptParams[bool]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			Render:       WrapRender[bool](&p, params.Render),
		}),
		Active:   params.Active,
		Inactive: params.Inactive,
	}

	actionHandler := NewActionHandler(map[Action]func(){
		UpAction:    p.toggleValue,
		DownAction:  p.toggleValue,
		LeftAction:  p.toggleValue,
		RightAction: p.toggleValue,
	}, nil)
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))
	})

	return &p
}

func (p *ConfirmPrompt) toggleValue() {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, 2)
	p.Value = !p.Value
}
