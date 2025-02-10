package core_test

import (
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPrompt() *core.SelectPrompt[string] {
	return core.NewSelectPrompt(core.SelectPromptParams[string]{
		Options: []*core.SelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
		Render: func(p *core.SelectPrompt[string]) string { return "" },
	})
}

func TestChangeSelectCursor(t *testing.T) {
	p := newSelectPrompt()

	assert.Equal(t, 0, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.EndKey})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.HomeKey})
	assert.Equal(t, 0, p.CursorIndex)

	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 0, p.CursorIndex)
}

func TestChangeSelectValue(t *testing.T) {
	p := newSelectPrompt()

	assert.Equal(t, p.Options[0].Value, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.Options[1].Value, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.Options[2].Value, p.Value)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, p.Options[1].Value, p.Value)
}

func TestSelectInitialValue(t *testing.T) {
	p := core.NewSelectPrompt(core.SelectPromptParams[string]{
		Options: []*core.SelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
		InitialValue: "baz",
		Render:       func(p *core.SelectPrompt[string]) string { return "" },
	})

	assert.Equal(t, "baz", p.Value)
	assert.Equal(t, 2, p.CursorIndex)
}

func TestLabelAsSelectValue(t *testing.T) {
	p := core.NewSelectPrompt(core.SelectPromptParams[string]{
		Options: []*core.SelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
		Render: func(p *core.SelectPrompt[string]) string { return "" },
	})

	assert.Equal(t, "foo", p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, "bar", p.Value)
}

func TestSelectRequiredValue(t *testing.T) {
	p := newSelectPrompt()
	p.Filter = true
	p.Required = true

	p.PressKey(&core.Key{Name: "#", Char: "#"})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}

func TestSelectFilter(t *testing.T) {
	p1 := newSelectPrompt()
	p2 := newSelectPrompt()
	p2.Filter = true

	p1.PressKey(&core.Key{Char: "b"})
	p2.PressKey(&core.Key{Char: "b"})

	assert.Greater(t, len(p1.Options), 0)
	assert.Greater(t, len(p2.Options), 0)
	assert.Greater(t, len(p1.Options), len(p2.Options))
}

func TestSelectFilterCursor(t *testing.T) {
	p := newSelectPrompt()
	p.Filter = true

	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Char: "z"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 2, p.CursorIndex)
}

func TestSelectFilterOutOptions(t *testing.T) {
	p := newSelectPrompt()
	p.Filter = true

	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Char: "z"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Char: "#"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 2, p.CursorIndex)
}
