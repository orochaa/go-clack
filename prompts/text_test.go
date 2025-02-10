package prompts_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/prompts"
	"github.com/orochaa/go-clack/prompts/symbols"
	"github.com/orochaa/go-clack/prompts/test"
	"github.com/stretchr/testify/assert"
)

const message = "test message"

func TestTextInitialState(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := symbols.BAR + " █"
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestTextInitialStateWithPlaceholder(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message, Placeholder: "foo"})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := symbols.BAR + " foo"
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestTextInitialStateWithInitialValue(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message, InitialValue: "foo"})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := symbols.BAR + " foo█"
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestTextErrorState(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message, InitialValue: "foo", Validate: func(value string) error {
		return fmt.Errorf("invalid value: %s", value)
	}})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := symbols.State(core.ErrorState) + " " + message
	valueWithCursor := symbols.BAR + " foo█"
	err := symbols.BAR_END + " invalid value: foo"
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, err}, "\r\n")
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestTextCancelState(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := symbols.State(core.CancelState) + " " + message
	value := symbols.BAR + " "
	expected := strings.Join([]string{symbols.BAR, title, value}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestTextCancelStateWithValue(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message, InitialValue: "foo"})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := symbols.State(core.CancelState) + " " + message
	value := symbols.BAR + " foo"
	expected := strings.Join([]string{symbols.BAR, title, value, symbols.BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestTextSubmitState(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message, InitialValue: "foo"})
	time.Sleep(time.Millisecond)

	p := test.TextTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := symbols.State(core.SubmitState) + " " + message
	value := symbols.BAR + " foo"
	expected := strings.Join([]string{symbols.BAR, title, value}, "\r\n")
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, expected, p.Frame)
}
