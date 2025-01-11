package prompts_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/stretchr/testify/assert"
)

func TestPathInitialState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := symbols.BAR + " " + p.ValueWithCursor()
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathInitialStateWithInitialValue(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message, InitialValue: "/foo"})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := symbols.BAR + " /foo█"
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathErrorState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message, InitialValue: "/foo", Validate: func(value string) error {
		return fmt.Errorf("invalid value: %s", value)
	}})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := symbols.State(core.ErrorState) + " " + message
	valueWithCursor := symbols.BAR + " /foo█"
	err := symbols.BAR_END + " invalid value: /foo"
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, err}, "\r\n")
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathCancelState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	p.Value = ""
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := symbols.State(core.CancelState) + " " + message
	value := symbols.BAR + " "
	expected := strings.Join([]string{symbols.BAR, title, value}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathCancelStateWithValue(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := symbols.State(core.CancelState) + " " + message
	value := symbols.BAR + " " + p.Value
	expected := strings.Join([]string{symbols.BAR, title, value, symbols.BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathSubmitState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := symbols.State(core.SubmitState) + " " + message
	value := symbols.BAR + " " + p.Value
	expected := strings.Join([]string{symbols.BAR, title, value}, "\r\n")
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathValueWithOptions(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.TabKey})
	p.PressKey(&core.Key{Name: core.TabKey})

	title := symbols.State(core.ActiveState) + " " + message
	value := symbols.BAR + " " + p.ValueWithCursor()
	options := p.FormatLines([]string{strings.Join(p.HintOptions, " ")}, core.FormatLinesOptions{
		Default: core.FormatLineOptions{
			Start: symbols.BAR,
		},
	})
	expected := strings.Join([]string{symbols.BAR, title, value, options, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.ActiveState, p.State)
	assert.Equal(t, expected, p.Frame)
}
