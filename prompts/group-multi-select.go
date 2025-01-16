package prompts

import (
	"context"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type GroupMultiSelectParams[TValue comparable] struct {
	Context        context.Context
	Message        string
	Options        map[string][]MultiSelectOption[TValue]
	InitialValue   []TValue
	DisabledGroups bool
	SpacedGroups   bool
	Required       bool
	Validate       func(value []TValue) error
}

// GroupMultiSelect displays a grouped multi select prompt to the user.
//
// The prompt displays a message within their options.
// The user can navigate between options using arrow keys.
// The user can select multiple options using space key.
// The prompt returns the values of the selected options.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context in which the prompt is displayed (default: nil).
//   - Message (string): The message to display to the user (default: "").
//   - Options (map[string][]MultiSelectOption[TValue]):
//     A map of group names to a slice of MultiSelectOption[TValue] values.
//     Each MultiSelectOption[TValue] contains a Label, Value, and IsSelected field (default: nil).
//   - InitialValue ([]TValue): The initial value of the prompt (default: nil).
//   - DisabledGroups (bool): Whether the groups selection are disabled (default: false).
//   - SpacedGroups (bool): Whether the groups are spaced out (default: false).
//   - Required (bool): Whether the prompt is required (default: false).
//   - Validate (func(value []TValue) error): Custom validation function for the prompt (default: nil).
//
// Returns:
//   - []TValue: The values of the selected options.
//   - error: An error if the user cancels the prompt or if an error occurs.
func GroupMultiSelect[TValue comparable](params GroupMultiSelectParams[TValue]) ([]TValue, error) {
	v := validator.NewValidator("GroupMultiSelect")
	v.ValidateOptions(len(params.Options))

	groups := make(map[string][]core.MultiSelectOption[TValue])
	for group, options := range params.Options {
		groups[group] = make([]core.MultiSelectOption[TValue], len(options))
		for i, option := range options {
			groups[group][i] = core.MultiSelectOption[TValue]{
				Label:      option.Label,
				Value:      option.Value,
				IsSelected: option.IsSelected,
			}
		}
	}

	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[TValue]{
		Context:        params.Context,
		InitialValue:   params.InitialValue,
		Options:        groups,
		DisabledGroups: params.DisabledGroups,
		Required:       params.Required,
		Validate:       params.Validate,
		Render: func(p *core.GroupMultiSelectPrompt[TValue]) string {
			var value string

			switch p.State {
			case core.SubmitState, core.CancelState:
				for _, option := range p.Options {
					if option.IsSelected {
						if value == "" {
							value = option.Label
						} else {
							value += ", " + option.Label
						}
					}
				}

			default:
				radioOptions := make([]string, len(p.Options))
				for i, option := range p.Options {
					if option.IsGroup {
						radioOptions[i] = groupOption(option, p.IsGroupSelected(option), i == p.CursorIndex, p.DisabledGroups)
						if params.SpacedGroups && i > 0 {
							radioOptions[i] = "\n" + radioOptions[i]
						}
						continue
					}

					radioOptions[i] = " " + groupOption(option, option.IsSelected, i == p.CursorIndex, false)
				}
				value = p.LimitLines(radioOptions, 3)
			}

			return theme.ApplyTheme(theme.ThemeParams[[]TValue]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.GroupMultiSelectTestingPrompt = p
	return p.Run()
}

func groupOption[TValue comparable](option *core.GroupMultiSelectOption[TValue], isSelected, isActive, isDisabled bool) string {
	var radio, label string

	if isSelected && isActive {
		radio = picocolors.Green(symbols.CHECKBOX_SELECTED)
		label = option.Label
	} else if isActive {
		radio = picocolors.Green(symbols.CHECKBOX_ACTIVE)
		label = option.Label
	} else if isSelected {
		radio = picocolors.Green(symbols.CHECKBOX_SELECTED)
		label = picocolors.Dim(option.Label)
	} else {
		radio = picocolors.Dim(symbols.CHECKBOX_INACTIVE)
		label = picocolors.Dim(option.Label)
	}

	if isDisabled {
		return label
	}

	return radio + " " + label
}
