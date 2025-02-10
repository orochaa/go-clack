package core_test

import (
	"fmt"
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/third_party/picocolors"

	"github.com/stretchr/testify/assert"
)

func newTextPrompt() *core.TextPrompt {
	return core.NewTextPrompt(core.TextPromptParams{
		Render: func(p *core.TextPrompt) string { return "" },
	})
}

func TestTextPromptInitialValue(t *testing.T) {
	p := core.NewTextPrompt(core.TextPromptParams{
		InitialValue: "foo",
		Render:       func(p *core.TextPrompt) string { return "" },
	})

	assert.Equal(t, "foo", p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTextPromptValueTrack(t *testing.T) {
	p := newTextPrompt()

	assert.Equal(t, "", p.Value)

	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)

	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)

	p.PressKey(&core.Key{Name: core.LeftKey})
	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, "b", p.Value)
}

func TestTextPromptValueWithCursor(t *testing.T) {
	p := newTextPrompt()
	cursor := "█"

	expected := cursor
	assert.Equal(t, expected, p.ValueWithCursor())

	p.Value = "foo"
	p.CursorIndex = len(p.Value)
	expected = p.Value + cursor
	assert.Equal(t, expected, p.ValueWithCursor())

	p.CursorIndex = len(p.Value) - 1
	expected = "fo" + picocolors.Inverse("o")
	assert.Equal(t, expected, p.ValueWithCursor())

	p.CursorIndex = 0
	expected = picocolors.Inverse("f") + "oo"
	assert.Equal(t, expected, p.ValueWithCursor())
}

func TestTextPromptValidate(t *testing.T) {
	p := core.NewTextPrompt(core.TextPromptParams{
		InitialValue: "123",
		Validate: func(value string) error {
			return fmt.Errorf("invalid value: %s", value)
		},
		Render: func(p *core.TextPrompt) string { return "" },
	})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, "invalid value: 123", p.Error)
}

func TestTextPromptPlaceholderCompletion(t *testing.T) {
	p := newTextPrompt()

	p.Placeholder = "foo"
	p.Value = ""
	p.PressKey(&core.Key{Name: core.TabKey})
	assert.Equal(t, "foo", p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Placeholder = "foo"
	p.Value = "bar"
	p.PressKey(&core.Key{Name: core.TabKey})
	assert.Equal(t, "bar", p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTextRequiredValue(t *testing.T) {
	p := newTextPrompt()
	p.Required = true

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}
