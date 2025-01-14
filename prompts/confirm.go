package prompts

import (
	"context"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type ConfirmParams struct {
	Context      context.Context
	Message      string
	InitialValue bool
	Active       string
	Inactive     string
}

func Confirm(params ConfirmParams) (bool, error) {
	p := core.NewConfirmPrompt(core.ConfirmPromptParams{
		Context:      params.Context,
		InitialValue: params.InitialValue,
		Active:       params.Active,
		Inactive:     params.Inactive,
		Render: func(p *core.ConfirmPrompt) string {
			activeRadio := picocolors.Green(symbols.RADIO_ACTIVE)
			inactiveRadio := picocolors.Dim(symbols.RADIO_INACTIVE)
			slash := picocolors.Dim("/")

			var value, valueWithCursor string
			if p.Value {
				value = p.Active
				valueWithCursor = strings.Join([]string{activeRadio, p.Active, slash, inactiveRadio, picocolors.Dim(p.Inactive)}, " ")
			} else {
				value = p.Inactive
				valueWithCursor = strings.Join([]string{inactiveRadio, picocolors.Dim(p.Active), slash, activeRadio, p.Inactive}, " ")
			}

			return theme.ApplyTheme(theme.ThemeParams[bool]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           value,
				ValueWithCursor: valueWithCursor,
			})
		},
	})
	test.ConfirmTestingPrompt = p
	return p.Run()
}
