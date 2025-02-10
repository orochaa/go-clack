package symbols

import (
	"github.com/orochaa/go-clack/core"
	isunicodesupported "github.com/orochaa/go-clack/third_party/is-unicode-supported"
)

func s(c, fallback string) string {
	if isunicodesupported.IsUnicodeSupported() {
		return c
	}
	return fallback
}

type Symbol = string

var (
	STEP_ACTIVE Symbol = s("◆", "*")
	STEP_CANCEL Symbol = s("■", "x")
	STEP_ERROR  Symbol = s("▲", "x")
	STEP_SUBMIT Symbol = s("◇", "o")

	BAR_START Symbol = s("┌", "T")
	BAR       Symbol = s("│", "|")
	BAR_END   Symbol = s("└", "—")

	RADIO_ACTIVE      Symbol = s("●", ">")
	RADIO_INACTIVE    Symbol = s("○", " ")
	CHECKBOX_ACTIVE   Symbol = s("◻", "[•]")
	CHECKBOX_SELECTED Symbol = s("◼", "[+]")
	CHECKBOX_INACTIVE Symbol = s("◻", "[ ]")
	PASSWORD_MASK     Symbol = s("▪", "•")

	BAR_H               Symbol = s("─", "-")
	CORNER_TOP_RIGHT    Symbol = s("╮", "+")
	CONNECT_LEFT        Symbol = s("├", "+")
	CORNER_BOTTOM_RIGHT Symbol = s("╯", "+")

	INFO    Symbol = s("●", "•")
	SUCCESS Symbol = s("◆", "*")
	WARN    Symbol = s("▲", "!")
	ERROR   Symbol = s("■", "x")
)

func State(state core.State) string {
	switch state {
	case core.ErrorState:
		return STEP_ERROR
	case core.CancelState:
		return STEP_CANCEL
	case core.SubmitState:
		return STEP_SUBMIT
	default:
		return STEP_ACTIVE
	}
}
