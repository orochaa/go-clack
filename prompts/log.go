package prompts

import (
	"fmt"
	"os"
	"strings"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/prompts/symbols"
	"github.com/orochaa/go-clack/third_party/picocolors"
)

type MessageLineOptions = core.FormatLineOptions

type LogOptions struct {
	Output *os.File
}

type MessageOptions struct {
	LogOptions
	core.FormatLinesOptions
}

func Message(msg string, options MessageOptions) {
	if options.Output == nil {
		options.Output = os.Stdout
	}
	p := &core.Prompt[string]{}
	formattedMsg := p.FormatLines(utils.SplitLines(msg), options.FormatLinesOptions)
	options.Output.WriteString(fmt.Sprintf("%s\r\n%s\r\n", picocolors.Gray(symbols.BAR), formattedMsg))
}

func styleMsg(msg string, style func(msg string) string) string {
	parts := utils.SplitLines(msg)
	styledParts := make([]string, len(parts))
	for i, part := range parts {
		styledParts[i] = style(part)
	}
	return strings.Join(styledParts, "\r\n")
}

// Intro displays an introductory message.
func Intro(msg string, options LogOptions) {
	p := &core.Prompt[string]{}
	formattedMsg := p.FormatLines(utils.SplitLines(msg), core.FormatLinesOptions{
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
func Cancel(msg string, options LogOptions) {
	Message(styleMsg(msg, picocolors.Red), MessageOptions{
		LogOptions: options,
		FormatLinesOptions: core.FormatLinesOptions{
			Default: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
			LastLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR_END),
			},
		},
	})
}

// Outro displays a closing message.
func Outro(msg string, options LogOptions) {
	Message("\r\n"+msg, MessageOptions{
		FormatLinesOptions: core.FormatLinesOptions{
			Default: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
			LastLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR_END),
			},
		},
	})
}

// Info displays an informational message with a blue info symbol.
func Info(msg string, options LogOptions) {
	Message(msg, MessageOptions{
		FormatLinesOptions: core.FormatLinesOptions{
			FirstLine: MessageLineOptions{
				Start: picocolors.Blue(symbols.INFO),
			},
			NewLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
		},
	})
}

// Success displays a success message with a green success symbol.
func Success(msg string, options LogOptions) {
	Message(msg, MessageOptions{
		FormatLinesOptions: core.FormatLinesOptions{
			FirstLine: MessageLineOptions{
				Start: picocolors.Green(symbols.SUCCESS),
			},
			NewLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
		},
	})
}

// Step displays a step message with a green step symbol.
func Step(msg string, options LogOptions) {
	Message(msg, MessageOptions{
		FormatLinesOptions: core.FormatLinesOptions{
			FirstLine: MessageLineOptions{
				Start: picocolors.Green(symbols.STEP_SUBMIT),
			},
			NewLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
		},
	})
}

// Warn displays a warning message with a yellow warning symbol.
func Warn(msg string, options LogOptions) {
	Message(msg, MessageOptions{
		FormatLinesOptions: core.FormatLinesOptions{
			FirstLine: MessageLineOptions{
				Start: picocolors.Yellow(symbols.WARN),
			},
			NewLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
		},
	})
}

// Error displays an error message with a red error symbol.
func Error(msg string, options LogOptions) {
	Message(msg, MessageOptions{
		FormatLinesOptions: core.FormatLinesOptions{
			FirstLine: MessageLineOptions{
				Start: picocolors.Red(symbols.ERROR),
			},
			NewLine: MessageLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
		},
	})
}
