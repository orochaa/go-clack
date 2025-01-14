package prompts

import (
	"context"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type PathParams struct {
	Context      context.Context
	Message      string
	InitialValue string
	OnlyShowDir  bool
	Required     bool
	Validate     func(value string) error
}

func Path(params PathParams) (string, error) {
	p := core.NewPathPrompt(core.PathPromptParams{
		Context:      params.Context,
		InitialValue: params.InitialValue,
		OnlyShowDir:  params.OnlyShowDir,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.PathPrompt) string {
			valueWithCursor := p.ValueWithCursor()

			if len(p.HintOptions) > 0 {
				var hintOptions string
				for i, hintOption := range p.HintOptions {
					if i == p.HintIndex {
						hintOptions += picocolors.Cyan(hintOption)
					} else {
						hintOptions += picocolors.Dim(hintOption)
					}
					if i+1 < len(p.HintOptions) {
						hintOptions += " "
					}
				}
				valueWithCursor += "\n" + hintOptions
			}

			return theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           p.Value,
				ValueWithCursor: valueWithCursor,
			})
		},
	})
	test.PathTestingPrompt = p
	return p.Run()
}
