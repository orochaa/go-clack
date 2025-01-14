package prompts

import (
	"context"
	"fmt"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type SelectKeyOption[TValue comparable] struct {
	Label string
	Value TValue
	Key   string
}

type SelectKeyParams[TValue comparable] struct {
	Context context.Context
	Message string
	Options []SelectKeyOption[TValue]
}

func SelectKey[TValue comparable](params SelectKeyParams[TValue]) (TValue, error) {
	v := validator.NewValidator("SelectKey")
	v.ValidateOptions(len(params.Options))

	var options []*core.SelectKeyOption[TValue]
	for _, option := range params.Options {
		options = append(options, &core.SelectKeyOption[TValue]{
			Label: option.Label,
			Value: option.Value,
			Key:   option.Key,
		})
	}

	p := core.NewSelectKeyPrompt(core.SelectKeyPromptParams[TValue]{
		Context: params.Context,
		Options: options,
		Render: func(p *core.SelectKeyPrompt[TValue]) string {
			var value string
			switch p.State {
			case core.SubmitState, core.CancelState:
			default:
				keyOptions := make([]string, len(params.Options))
				for i, option := range params.Options {
					key := picocolors.Cyan("[" + option.Key + "]")
					label := option.Label
					keyOptions[i] = fmt.Sprintf("%s %s", key, label)
				}
				value = p.LimitLines(keyOptions, 3)
			}

			return theme.ApplyTheme(theme.ThemeParams[TValue]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           params.Options[p.CursorIndex].Label,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectKeyTestingPrompt = p
	return p.Run()
}
