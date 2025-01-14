package prompts

import (
	"context"
	"fmt"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
	Hint  string
}

type SelectParams[TValue comparable] struct {
	Context      context.Context
	Message      string
	InitialValue TValue
	Options      []*SelectOption[TValue]
	Filter       bool
	Required     bool
}

func Select[TValue comparable](params SelectParams[TValue]) (TValue, error) {
	v := validator.NewValidator("Select")
	v.ValidateOptions(len(params.Options))

	var options []*core.SelectOption[TValue]
	for _, option := range params.Options {
		options = append(options, &core.SelectOption[TValue]{
			Label: option.Label,
			Value: option.Value,
		})
	}

	p := core.NewSelectPrompt(core.SelectPromptParams[TValue]{
		Context:      params.Context,
		InitialValue: params.InitialValue,
		Options:      options,
		Filter:       params.Filter,
		Required:     params.Required,
		Render: func(p *core.SelectPrompt[TValue]) string {
			message := params.Message
			var value string

			switch p.State {
			case core.SubmitState, core.CancelState:
				if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
					value = p.Options[p.CursorIndex].Label
				}
			default:
				radioOptions := make([]string, len(p.Options))
				for _, option := range params.Options {
					for i, _option := range p.Options {
						if option.Label != _option.Label {
							continue
						}

						if i == p.CursorIndex && option.Hint != "" {
							radio := picocolors.Green(symbols.RADIO_ACTIVE)
							label := option.Label
							hint := picocolors.Dim("(" + option.Hint + ")")
							radioOptions[i] = fmt.Sprintf("%s %s %s", radio, label, hint)
						} else if i == p.CursorIndex {
							radio := picocolors.Green(symbols.RADIO_ACTIVE)
							label := option.Label
							radioOptions[i] = fmt.Sprintf("%s %s", radio, label)
						} else {
							radio := picocolors.Dim(symbols.RADIO_INACTIVE)
							label := picocolors.Dim(option.Label)
							radioOptions[i] = fmt.Sprintf("%s %s", radio, label)
						}

						break
					}
				}

				if p.Filter {
					if p.Search == "" {
						message = fmt.Sprintf("%s\n> %s", message, picocolors.Inverse("T")+picocolors.Dim("ype to filter..."))
					} else {
						message = fmt.Sprintf("%s\n> %s", message, p.Search+"█")
					}

					value = p.LimitLines(radioOptions, 4)
					break
				}

				value = p.LimitLines(radioOptions, 3)
			}

			return theme.ApplyTheme(theme.ThemeParams[TValue]{
				Ctx:             p.Prompt,
				Message:         message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectTestingPrompt = p
	return p.Run()
}
