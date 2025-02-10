package theme

import (
	"strings"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/prompts/symbols"
	"github.com/orochaa/go-clack/third_party/picocolors"
)

type ThemeValue interface {
	string | any | []any
}

type ThemeParams[TValue ThemeValue] struct {
	Context         core.Prompt[TValue]
	Message         string
	Value           string
	ValueWithCursor string
	Placeholder     string
}

func ApplyTheme[TValue ThemeValue](params ThemeParams[TValue]) string {
	ctx := params.Context

	frame := make([]string, 0, 4)
	frame = append(frame, picocolors.Gray(symbols.BAR))

	symbolColor := SymbolColor(ctx.State)
	barColor := BarColor(ctx.State)
	title := ctx.FormatLines(utils.SplitLines(params.Message), core.FormatLinesOptions{
		FirstLine: core.FormatLineOptions{
			Start: symbolColor(symbols.State(ctx.State)),
		},
		NewLine: core.FormatLineOptions{
			Start: barColor(symbols.BAR),
		},
	})
	frame = append(frame, title)

	var valueWithCursor string
	if params.Placeholder != "" && (params.ValueWithCursor == "" || (ctx.State == core.InitialState && params.ValueWithCursor == "█")) {
		valueWithCursor = picocolors.Inverse(string(params.Placeholder[0])) + picocolors.Dim(params.Placeholder[1:])
	} else {
		valueWithCursor = params.ValueWithCursor
	}

	switch ctx.State {
	case core.ErrorState:
		value := ctx.FormatLines(utils.SplitLines(valueWithCursor), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
			},
		})
		frame = append(frame, value)

		if ctx.Error != "" {
			err := ctx.FormatLines(utils.SplitLines(ctx.Error), core.FormatLinesOptions{
				Default: core.FormatLineOptions{
					Start: barColor(symbols.BAR),
					Style: picocolors.Yellow,
				},
				LastLine: core.FormatLineOptions{
					Start: barColor(symbols.BAR_END),
				},
			})
			frame = append(frame, err)
		}

	case core.CancelState:
		value := ctx.FormatLines(utils.SplitLines(params.Value), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: func(line string) string {
					return picocolors.Strikethrough(picocolors.Dim(line))
				},
			},
		})
		frame = append(frame, value)

		if params.Value != "" {
			end := barColor(symbols.BAR)
			frame = append(frame, end)
		}

	case core.SubmitState:
		value := ctx.FormatLines(utils.SplitLines(params.Value), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: picocolors.Dim,
			},
		})
		frame = append(frame, value)

	case core.ValidateState:
		value := ctx.FormatLines(utils.SplitLines(params.Value), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: picocolors.Dim,
			},
		})
		dots := strings.Repeat(".", int(ctx.ValidationDuration.Seconds())%4)
		validatingMsg := barColor(symbols.BAR_END) + " " + picocolors.Dim("validating"+dots)
		frame = append(frame, value, validatingMsg)

	default:
		value := ctx.FormatLines(utils.SplitLines(valueWithCursor), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
			},
		})
		end := barColor(symbols.BAR_END)
		frame = append(frame, value, end)
	}

	return strings.Join(frame, "\r\n")
}

func SymbolColor(state core.State) func(input string) string {
	switch state {
	case core.ErrorState, core.CancelState:
		return picocolors.Red
	case core.SubmitState:
		return picocolors.Green
	default:
		return picocolors.Cyan
	}
}

func BarColor(state core.State) func(input string) string {
	switch state {
	case core.ErrorState:
		return picocolors.Yellow
	case core.InitialState, core.ActiveState:
		return picocolors.Cyan
	default:
		return picocolors.Gray
	}
}
