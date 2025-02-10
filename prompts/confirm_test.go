package prompts_test

import (
	"strings"
	"testing"
	"time"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/prompts"
	"github.com/orochaa/go-clack/prompts/symbols"
	"github.com/orochaa/go-clack/prompts/test"
	"github.com/stretchr/testify/assert"
)

func TestConfirmInitialState(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := strings.Join([]string{symbols.BAR, symbols.RADIO_INACTIVE, p.Active, "/", symbols.RADIO_ACTIVE, p.Inactive}, " ")
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmInitialStateWithInitialValue(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message, InitialValue: true})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	title := symbols.State(core.InitialState) + " " + message
	valueWithCursor := strings.Join([]string{symbols.BAR, symbols.RADIO_ACTIVE, p.Active, "/", symbols.RADIO_INACTIVE, p.Inactive}, " ")
	expected := strings.Join([]string{symbols.BAR, title, valueWithCursor, symbols.BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmCancelState(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := symbols.State(core.CancelState) + " " + message
	value := symbols.BAR + " " + p.Inactive
	expected := strings.Join([]string{symbols.BAR, title, value, symbols.BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmCancelStateWithValue(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message, InitialValue: true})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := symbols.State(core.CancelState) + " " + message
	value := symbols.BAR + " " + p.Active
	expected := strings.Join([]string{symbols.BAR, title, value, symbols.BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmSubmitState(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := symbols.State(core.SubmitState) + " " + message
	value := symbols.BAR + " " + p.Inactive
	expected := strings.Join([]string{symbols.BAR, title, value}, "\r\n")
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, expected, p.Frame)
}
