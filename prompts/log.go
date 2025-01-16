package prompts

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type MessageLineOptions = core.FormatLineOptions
type MessageOptions = core.FormatLinesOptions

func Message(msg string, options MessageOptions) {
	p := &core.Prompt[string]{}
	formattedMsg := p.FormatLines(strings.Split(msg, "\n"), options)
	os.Stdout.WriteString(fmt.Sprintf("%s\r\n%s\r\n", picocolors.Gray(symbols.BAR), formattedMsg))
}

func styleMsg(msg string, style func(msg string) string) string {
	parts := strings.Split(msg, "\n")
	styledParts := make([]string, len(parts))
	for i, part := range parts {
		styledParts[i] = style(part)
	}
	return strings.Join(styledParts, "\n")
}

// Intro displays an introductory message.
func Intro(msg string) {
	p := &core.Prompt[string]{}
	formattedMsg := p.FormatLines(strings.Split(msg, "\n"), MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR_START),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
	})
	os.Stdout.WriteString(fmt.Sprintf("\r\n%s\r\n%s\r\n", formattedMsg, picocolors.Gray(symbols.BAR)))
}

// Cancel displays a cancellation message styled in red.
func Cancel(msg string) {
	Message(styleMsg(msg, picocolors.Red), MessageOptions{
		Default: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
		LastLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR_END),
		},
	})
}

// Outro displays a closing message.
func Outro(msg string) {
	Message("\n"+msg, MessageOptions{
		Default: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
		LastLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR_END),
		},
	})
}

// Info displays an informational message with a blue info symbol.
func Info(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Blue(symbols.INFO),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
	})
}

// Success displays a success message with a green success symbol.
func Success(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Green(symbols.SUCCESS),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
	})
}

// Step displays a step message with a green step symbol.
func Step(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Green(symbols.STEP_SUBMIT),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
	})
}

// Warn displays a warning message with a yellow warning symbol.
func Warn(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Yellow(symbols.WARN),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
	})
}

// Error displays an error message with a red error symbol.
func Error(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Red(symbols.ERROR),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(symbols.BAR),
		},
	})
}
