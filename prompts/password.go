package prompts

import (
	"context"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
)

type PasswordParams struct {
	Context      context.Context
	Message      string
	InitialValue string
	Required     bool
	Validate     func(value string) error
}

func Password(params PasswordParams) (string, error) {
	p := core.NewPasswordPrompt(core.PasswordPromptParams{
		Context:      params.Context,
		InitialValue: params.InitialValue,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.PasswordPrompt) string {
			return theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           p.ValueWithMask(),
				ValueWithCursor: p.ValueWithMaskAndCursor(),
			})
		},
	})
	test.PasswordTestingPrompt = p
	return p.Run()
}
