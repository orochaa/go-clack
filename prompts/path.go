package prompts

import (
	"context"
	"os"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type PathParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	Message      string
	InitialValue string
	OnlyShowDir  bool
	Required     bool
	Validate     func(value string) error
}

// Path displays a input prompt to the user.
//
// The prompt displays a message.
// The user can input a path.
// The prompt has built-in autosuggestion and autocomplete features.
// The prompt returns the path.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Message (string): The message to display to the user (default: "").
//   - InitialValue (string): The initial value of the path input (default: current working directory).
//   - OnlyShowDir (bool): Whether to only show directories (default: false).
//   - Required (bool): Whether the path input is required (default: false).
//   - FileSystem (FileSystem): The file system implementation to use (default: OSFileSystem).
//   - Validate (func(value string) error): Custom validation function for the path (default: nil).
//
// Returns:
//   - string: The path value.
//   - error: An error if the user cancels the prompt or if an error occurs.
func Path(params PathParams) (string, error) {
	p := core.NewPathPrompt(core.PathPromptParams{
		Context:      params.Context,
		Input:        params.Input,
		Output:       params.Output,
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
				valueWithCursor += "\r\n" + hintOptions
			}

			return theme.ApplyTheme(theme.ThemeParams[string]{
				Context:         p.Prompt,
				Message:         params.Message,
				Value:           p.Value,
				ValueWithCursor: valueWithCursor,
			})
		},
	})
	test.PathTestingPrompt = p
	return p.Run()
}
