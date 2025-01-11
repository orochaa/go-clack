package theme_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/stretchr/testify/assert"
)

func TestApplyThemeInitialState(t *testing.T) {
	testCases := []struct {
		Description     string
		ValueWithCursor string
		Placeholder     string
		Expected        string
	}{
		{
			Description: "InitialState",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description: "InitialStateWithPlaceholder",
			Placeholder: "Placeholder",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Placeholder",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description:     "InitialStateWithValue",
			ValueWithCursor: "Value with cursor█",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Value with cursor█",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description:     "InitialStateWithPlaceholderAndValue",
			Placeholder:     "Placeholder",
			ValueWithCursor: "Value with cursor█",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Value with cursor█",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description:     "InitialStateWithPlaceholderAndCursor",
			Placeholder:     "Placeholder",
			ValueWithCursor: "█",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Placeholder",
				symbols.BAR_END,
			}, "\r\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Description, func(t *testing.T) {
			frame := theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx: core.Prompt[string]{
					State: core.InitialState,
				},
				Message:         "Test message",
				ValueWithCursor: tC.ValueWithCursor,
				Placeholder:     tC.Placeholder,
			})
			assert.Equal(t, tC.Expected, frame)
		})
	}
}

func TestApplyThemeActiveState(t *testing.T) {
	testCases := []struct {
		Description     string
		ValueWithCursor string
		Placeholder     string
		Expected        string
	}{
		{
			Description: "ActiveState",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description: "ActiveStateWithPlaceholder",
			Placeholder: "Placeholder",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Placeholder",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description:     "ActiveStateWithValue",
			ValueWithCursor: "Value with cursor ",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Value with cursor ",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description:     "ActiveStateWithPlaceholderAndValue",
			Placeholder:     "Placeholder",
			ValueWithCursor: "Value with cursor ",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Value with cursor ",
				symbols.BAR_END,
			}, "\r\n"),
		},
		{
			Description:     "ActiveStateWithPlaceholderAndCursor",
			Placeholder:     "Placeholder",
			ValueWithCursor: " ",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END,
			}, "\r\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Description, func(t *testing.T) {
			frame := theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx: core.Prompt[string]{
					State: core.ActiveState,
				},
				Message:         "Test message",
				ValueWithCursor: tC.ValueWithCursor,
				Placeholder:     tC.Placeholder,
			})
			assert.Equal(t, tC.Expected, frame)
		})
	}
}

func TestApplyThemeErrorState(t *testing.T) {
	testCases := []struct {
		Description     string
		ValueWithCursor string
		Placeholder     string
		Error           string
		Expected        string
	}{
		{
			Description: "ErrorState",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " ",
			}, "\r\n"),
		},
		{
			Description: "ErrorStateWithPlaceholder",
			Placeholder: "Placeholder",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " Placeholder",
			}, "\r\n"),
		},
		{
			Description: "ErrorStateWithPlaceholderAndError",
			Placeholder: "Placeholder",
			Error:       "Error message",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " Placeholder",
				symbols.BAR_END + " Error message",
			}, "\r\n"),
		},
		{
			Description:     "ErrorStateWithValue",
			ValueWithCursor: "Value with cursor ",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " Value with cursor ",
			}, "\r\n"),
		},
		{
			Description:     "ErrorStateWithValueAndError",
			ValueWithCursor: "Value with cursor ",
			Error:           "Error message",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " Value with cursor ",
				symbols.BAR_END + " Error message",
			}, "\r\n"),
		},
		{
			Description:     "ErrorStateWithPlaceholderAndValue",
			Placeholder:     "Placeholder",
			ValueWithCursor: "Value with cursor ",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " Value with cursor ",
			}, "\r\n"),
		},
		{
			Description:     "ErrorStateWithWithPlaceholderAndValueAndError",
			Placeholder:     "Placeholder",
			ValueWithCursor: "Value with cursor ",
			Error:           "Error message",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ERROR + " Test message",
				symbols.BAR + " Value with cursor ",
				symbols.BAR_END + " Error message",
			}, "\r\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Description, func(t *testing.T) {
			frame := theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx: core.Prompt[string]{
					State: core.ErrorState,
					Error: tC.Error,
				},
				Message:         "Test message",
				ValueWithCursor: tC.ValueWithCursor,
				Placeholder:     tC.Placeholder,
			})
			assert.Equal(t, tC.Expected, frame)
		})
	}
}

func TestApplyThemeSubmitState(t *testing.T) {
	testCases := []struct {
		Description string
		Value       string
		Placeholder string
		Expected    string
	}{
		{
			Description: "SubmitState",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_SUBMIT + " Test message",
				symbols.BAR + " ",
			}, "\r\n"),
		},
		{
			Description: "SubmitStateWithPlaceholder",
			Placeholder: "Placeholder",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_SUBMIT + " Test message",
				symbols.BAR + " ",
			}, "\r\n"),
		},
		{
			Description: "SubmitStateWithValue",
			Value:       "Value",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_SUBMIT + " Test message",
				symbols.BAR + " Value",
			}, "\r\n"),
		},
		{
			Description: "SubmitStateWithPlaceholderAndValue",
			Placeholder: "Placeholder",
			Value:       "Value",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_SUBMIT + " Test message",
				symbols.BAR + " Value",
			}, "\r\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Description, func(t *testing.T) {
			frame := theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx: core.Prompt[string]{
					State: core.SubmitState,
				},
				Message:     "Test message",
				Value:       tC.Value,
				Placeholder: tC.Placeholder,
			})
			assert.Equal(t, tC.Expected, frame)
		})
	}
}

func TestApplyThemeCancelState(t *testing.T) {
	testCases := []struct {
		Description string
		Value       string
		Placeholder string
		Expected    string
	}{
		{
			Description: "CancelState",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_CANCEL + " Test message",
				symbols.BAR + " ",
			}, "\r\n"),
		},
		{
			Description: "CancelStateWithPlaceholder",
			Placeholder: "Placeholder",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_CANCEL + " Test message",
				symbols.BAR + " ",
			}, "\r\n"),
		},
		{
			Description: "CancelStateWithValue",
			Value:       "Value",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_CANCEL + " Test message",
				symbols.BAR + " Value",
				symbols.BAR,
			}, "\r\n"),
		},
		{
			Description: "CancelStateWithPlaceholderAndValue",
			Placeholder: "Placeholder",
			Value:       "Value",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_CANCEL + " Test message",
				symbols.BAR + " Value",
				symbols.BAR,
			}, "\r\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Description, func(t *testing.T) {
			frame := theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx: core.Prompt[string]{
					State: core.CancelState,
				},
				Message:     "Test message",
				Value:       tC.Value,
				Placeholder: tC.Placeholder,
			})
			assert.Equal(t, tC.Expected, frame)
		})
	}
}

func TestApplyThemeValidateState(t *testing.T) {
	testCases := []struct {
		Description        string
		Value              string
		Placeholder        string
		ValidationDuration time.Duration
		Expected           string
	}{
		{
			Description: "ValidateState",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating",
			}, "\r\n"),
		},
		{
			Description:        "ValidateStateAfterOneSecond",
			ValidationDuration: 1 * time.Second,
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating.",
			}, "\r\n"),
		},
		{
			Description:        "ValidateStateAfterTwoSecond",
			ValidationDuration: 2 * time.Second,
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating..",
			}, "\r\n"),
		},
		{
			Description:        "ValidateStateAfterThreeSecond",
			ValidationDuration: 3 * time.Second,
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating...",
			}, "\r\n"),
		},
		{
			Description:        "ValidateStateAfterFourSecond",
			ValidationDuration: 4 * time.Second,
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating",
			}, "\r\n"),
		},
		{
			Description:        "ValidateStateAfterFiveSecond",
			ValidationDuration: 5 * time.Second,
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating.",
			}, "\r\n"),
		},
		{
			Description: "ValidateStateWithPlaceholder",
			Placeholder: "Placeholder",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " ",
				symbols.BAR_END + " validating",
			}, "\r\n"),
		},
		{
			Description: "ValidateStateWithValue",
			Value:       "Value",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Value",
				symbols.BAR_END + " validating",
			}, "\r\n"),
		},
		{
			Description: "ValidateStateWithPlaceholderAndValue",
			Placeholder: "Placeholder",
			Value:       "Value",
			Expected: strings.Join([]string{
				symbols.BAR,
				symbols.STEP_ACTIVE + " Test message",
				symbols.BAR + " Value",
				symbols.BAR_END + " validating",
			}, "\r\n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Description, func(t *testing.T) {
			frame := theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx: core.Prompt[string]{
					State:              core.ValidateState,
					IsValidating:       true,
					ValidationDuration: tC.ValidationDuration,
				},
				Message:     "Test message",
				Value:       tC.Value,
				Placeholder: tC.Placeholder,
			})
			assert.Equal(t, tC.Expected, frame)
		})
	}
}
