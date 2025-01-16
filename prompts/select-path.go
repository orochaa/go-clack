package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type FileSystem = core.FileSystem

type SelectPathParams struct {
	Context      context.Context
	Message      string
	InitialValue string
	OnlyShowDir  bool
	Filter       bool
	FileSystem   FileSystem
}

// SelectPath displays a select prompt to the user.
//
// The prompt displays a message within their options.
// The user can navigate through directories and files using arrow keys.
// The user can select an option using enter key.
// The prompt returns the path of the selected option.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context in which the prompt is displayed (default: nil).
//   - Message (string): The message to display to the user (default: "").
//   - InitialValue (string): The initial path value (default: current working directory).
//   - OnlyShowDir (bool): Whether to only show directories (default: false).
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - FileSystem (FileSystem): The file system implementation to use (default: OSFileSystem).
//
// Returns:
//   - string: The path of the selected option.
//   - error: An error if the user cancels the prompt or if an error occurs.
func SelectPath(params SelectPathParams) (string, error) {
	p := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Context:      params.Context,
		InitialValue: params.InitialValue,
		OnlyShowDir:  params.OnlyShowDir,
		Filter:       params.Filter,
		FileSystem:   params.FileSystem,
		Render: func(p *core.SelectPathPrompt) string {
			message := params.Message
			var value string

			switch p.State {
			case core.SubmitState, core.CancelState:
			default:
				options := p.Options()
				radioOptions := make([]string, len(options))
				for i, option := range options {
					var radio, label, dir string
					if option.IsDir && option.IsOpen {
						dir = "v"
					} else if option.IsDir {
						dir = ">"
					}
					if option.IsEqual(p.CurrentOption) {
						radio = picocolors.Green(symbols.RADIO_ACTIVE)
						label = option.Name
					} else {
						radio = picocolors.Dim(symbols.RADIO_INACTIVE)
						label = picocolors.Dim(option.Name)
						dir = picocolors.Dim(dir)
					}
					depth := strings.Repeat(" ", option.Depth)
					radioOptions[i] = fmt.Sprintf("%s%s %s %s", depth, radio, label, dir)
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

			return theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         message,
				Value:           p.Value,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectPathTestingPrompt = p
	return p.Run()
}
