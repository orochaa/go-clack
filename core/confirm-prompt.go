package core

import (
	"context"
	"os"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/core/validator"
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

// NewConfirmPrompt initializes and returns a new instance of ConfirmPrompt.
//
// The user can toggle between the two options using arrow keys.
// The prompt returns the selected value if the user confirms their choice.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Active (string): The label displayed when the prompt is in the "active" (true) state (default: "yes").
//   - Inactive (string): The label displayed when the prompt is in the "inactive" (false) state (default: "no").
//   - InitialValue (bool): The initial value of the prompt (default: false).
//   - Render (func(p *MultiSelectPathPrompt) string): A custom render function for the prompt (default: nil).
//
// Returns:
//   - *ConfirmPrompt: A pointer to the newly created ConfirmPrompt instance.
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

// toggleValue toggles the current value of the ConfirmPrompt between true and false.
// It also updates the cursor index to ensure it stays within valid bounds.
func (p *ConfirmPrompt) toggleValue() {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, 2)
	p.Value = !p.Value
}
