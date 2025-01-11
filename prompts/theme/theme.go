package theme

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type ThemeValue interface {
	string | any | []any
}

type ThemeParams[TValue ThemeValue] struct {
	Ctx             core.Prompt[TValue]
	Message         string
	Value           string
	ValueWithCursor string
	Placeholder     string
}

func ApplyTheme[TValue ThemeValue](params ThemeParams[TValue]) string {
	ctx := params.Ctx

	symbolColor := SymbolColor(ctx.State)
	barColor := BarColor(ctx.State)

	title := strings.Join([]string{
		picocolors.Gray(symbols.BAR),
		ctx.FormatLines(strings.Split(params.Message, "\n"), core.FormatLinesOptions{
			FirstLine: core.FormatLineOptions{
				Start: symbolColor(symbols.State(ctx.State)),
			},
			NewLine: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
			},
		}),
	}, "\r\n")

	var valueWithCursor string
	if params.Placeholder != "" && (params.ValueWithCursor == "" || (ctx.State == core.InitialState && params.ValueWithCursor == "█")) {
		valueWithCursor = picocolors.Inverse(string(params.Placeholder[0])) + picocolors.Dim(params.Placeholder[1:])
	} else {
		valueWithCursor = params.ValueWithCursor
	}

	switch ctx.State {
	case core.ErrorState:
		value := ctx.FormatLines(strings.Split(valueWithCursor, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
			},
		})
		if ctx.Error == "" {
			return strings.Join([]string{title, value}, "\r\n")
		}
		err := ctx.FormatLines(strings.Split(ctx.Error, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: picocolors.Yellow,
			},
			LastLine: core.FormatLineOptions{
				Start: barColor(symbols.BAR_END),
			},
		})
		return strings.Join([]string{title, value, err}, "\r\n")

	case core.CancelState:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: func(line string) string {
					return picocolors.Strikethrough(picocolors.Dim(line))
				},
			},
		})
		if params.Value == "" {
			return strings.Join([]string{title, value}, "\r\n")
		}
		end := barColor(symbols.BAR)
		return strings.Join([]string{title, value, end}, "\r\n")

	case core.SubmitState:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: picocolors.Dim,
			},
		})
		return strings.Join([]string{title, value}, "\r\n")

	case core.ValidateState:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
				Style: picocolors.Dim,
			},
		})
		dots := strings.Repeat(".", int(ctx.ValidationDuration.Seconds())%4)
		validatingMsg := barColor(symbols.BAR_END) + " " + picocolors.Dim("validating"+dots)
		return strings.Join([]string{title, value, validatingMsg}, "\r\n")

	default:
		value := ctx.FormatLines(strings.Split(valueWithCursor, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: barColor(symbols.BAR),
			},
		})
		end := barColor(symbols.BAR_END)

		return strings.Join([]string{title, value, end}, "\r\n")
	}
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
