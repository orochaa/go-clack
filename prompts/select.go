package prompts

import (
	"context"
	"fmt"
	"os"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/core/validator"
	"github.com/orochaa/go-clack/prompts/symbols"
	"github.com/orochaa/go-clack/prompts/test"
	"github.com/orochaa/go-clack/prompts/theme"
	"github.com/orochaa/go-clack/third_party/picocolors"
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
	Hint  string
}

type SelectParams[TValue comparable] struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	Message      string
	InitialValue TValue
	Options      []*SelectOption[TValue]
	Filter       bool
	Required     bool
}

// Select displays a select prompt to the user.
//
// The prompt displays a message within their options.
// The user can navigate through options using arrow keys.
// The user can select an option using enter key.
// The prompt returns the value of the selected option.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Message (string): The message to display to the user (default: "").
//   - InitialValue (TValue): The initial value of the prompt (default: zero value of TValue).
//   - Options ([]*SelectOption[TValue]): A list of options for the prompt (default: nil).
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - Required (bool): Whether the prompt requires a selection (default: false).
//
// Returns:
//   - TValue: The value of the selected option.
//   - error: An error if the user cancels the prompt or if an error occurs.
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
		Input:        params.Input,
		Output:       params.Output,
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
				Context:         p.Prompt,
				Message:         message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectTestingPrompt = p
	return p.Run()
}
