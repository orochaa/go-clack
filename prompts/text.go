package prompts

import (
	"context"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
)

type TextParams struct {
	Context      context.Context
	Message      string
	Placeholder  string
	InitialValue string
	Required     bool
	Validate     func(value string) error
}

func Text(params TextParams) (string, error) {
	p := core.NewTextPrompt(core.TextPromptParams{
		Context:      params.Context,
		InitialValue: params.InitialValue,
		Placeholder:  params.Placeholder,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.TextPrompt) string {
			return theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx:             p.Prompt,
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
