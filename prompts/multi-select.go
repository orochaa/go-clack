package prompts

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type MultiSelectOption[TValue comparable] struct {
	Label      string
	Value      TValue
	IsSelected bool
	Hint       string
}

type MultiSelectParams[TValue comparable] struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	Message      string
	Options      []*MultiSelectOption[TValue]
	InitialValue []TValue
	Filter       bool
	Required     bool
	Validate     func(value []TValue) error
}

// MultiSelect displays a multi-select prompt to the user.
//
// The prompt displays a message within their options.
// The user can navigate through options using arrow keys.
// The user can select multiple options using space key.
// The prompt returns the value of the selected options.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Message (string): The message to display to the user (default: "").
//   - Options ([]*MultiSelectOption[TValue]): A list of options for the prompt (default: nil.
//   - InitialValue ([]TValue): The initial selected values (default: nil.
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - Required (bool): Whether the prompt requires at least one selection (default: false).
//   - Validate (func(value []TValue) error): Custom validation function for the prompt (default: nil).
//
// Returns:
//   - []TValue: A slice of values of the selected options.
//   - error: An error if the user cancels the prompt or if an error occurs.
func MultiSelect[TValue comparable](params MultiSelectParams[TValue]) ([]TValue, error) {
	v := validator.NewValidator("MultiSelect")
	v.ValidateOptions(len(params.Options))

	var options []*core.MultiSelectOption[TValue]
	for _, option := range params.Options {
		options = append(options, &core.MultiSelectOption[TValue]{
			Label:      option.Label,
			Value:      option.Value,
			IsSelected: option.IsSelected,
		})
	}

	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[TValue]{
		Context:      params.Context,
		Input:        params.Input,
		Output:       params.Output,
		InitialValue: params.InitialValue,
		Options:      options,
		Filter:       params.Filter,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.MultiSelectPrompt[TValue]) string {
			message := params.Message
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
					var radio, label, hint string
					if option.IsSelected && i == p.CursorIndex {
						radio = picocolors.Green(symbols.CHECKBOX_SELECTED)
						label = option.Label
						if params.Options[i].Hint != "" {
							hint = picocolors.Dim("(" + params.Options[i].Hint + ")")
						}
					} else if i == p.CursorIndex {
						radio = picocolors.Green(symbols.CHECKBOX_ACTIVE)
						label = option.Label
						if params.Options[i].Hint != "" {
							hint = picocolors.Dim("(" + params.Options[i].Hint + ")")
						}
					} else if option.IsSelected {
						radio = picocolors.Green(symbols.CHECKBOX_SELECTED)
						label = picocolors.Dim(option.Label)
					} else {
						radio = picocolors.Dim(symbols.CHECKBOX_INACTIVE)
						label = picocolors.Dim(option.Label)
					}
					radioOptions[i] = strings.Join([]string{radio, label, hint}, " ")
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

			return theme.ApplyTheme(theme.ThemeParams[[]TValue]{
				Ctx:             p.Prompt,
				Message:         message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.MultiSelectTestingPrompt = p
	return p.Run()
}
