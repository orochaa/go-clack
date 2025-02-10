package core_test

import (
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPrompt() *core.MultiSelectPrompt[string] {
	return core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
		Render: func(p *core.MultiSelectPrompt[string]) string { return "" },
	})
}

func TestChangeMultiSelectCursor(t *testing.T) {
	p := newMultiSelectPrompt()

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

func TestChangeMultiSelectValue(t *testing.T) {
	p := newMultiSelectPrompt()

	assert.Equal(t, []string(nil), p.Value)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[0].Value}, p.Value)
	assert.Equal(t, true, p.Options[0].IsSelected)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{}, p.Value)
	assert.Equal(t, false, p.Options[0].IsSelected)

	expected := make([]string, len(p.Options))
	for i, option := range p.Options {
		expected[i] = option.Value
	}
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, expected, p.Value)
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, []string{}, p.Value)

	for _, option := range p.Options {
		option.IsSelected = true
		p.Value = append(p.Value, option.Value)
	}
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: core.SpaceKey})
	expected = append([]string{expected[0]}, expected[2:]...)
	assert.Equal(t, expected, p.Value)
}

func TestMultiSelectInitialValue(t *testing.T) {
	initialValue := []string{"foo", "baz"}
	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
			{Value: "foo"},
			{Value: "bar"},
			{Value: "baz"},
		},
		InitialValue: initialValue,
		Render:       func(p *core.MultiSelectPrompt[string]) string { return "" },
	})

	assert.Equal(t, initialValue, p.Value)

	for _, value := range p.Value {
		for _, option := range p.Options {
			if option.Value == value {
				assert.True(t, option.IsSelected, option.Value)
			}
		}
	}
}

func TestMultiSelectInitialSelectedOptions(t *testing.T) {
	initialValue := []string{"foo", "baz"}
	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
			{Value: "foo", IsSelected: true},
			{Value: "bar"},
			{Value: "baz", IsSelected: true},
		},
		Render: func(p *core.MultiSelectPrompt[string]) string { return "" },
	})

	assert.Equal(t, initialValue, p.Value)

	for _, value := range p.Value {
		for _, option := range p.Options {
			if option.Value == value {
				assert.True(t, option.IsSelected, option.Value)
			}
		}
	}
}

func TestLabelAsMultiSelectValue(t *testing.T) {
	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
		Render: func(p *core.MultiSelectPrompt[string]) string { return "" },
	})

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{"foo"}, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{"foo", "bar"}, p.Value)
}

func TestMultiSelectRequiredValue(t *testing.T) {
	p := newMultiSelectPrompt()
	p.Required = true

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}

func TestMultiSelectFilter(t *testing.T) {
	p1 := newMultiSelectPrompt()
	p2 := newMultiSelectPrompt()
	p2.Filter = true

	p1.PressKey(&core.Key{Char: "b"})
	p2.PressKey(&core.Key{Char: "b"})

	assert.Greater(t, len(p1.Options), 0)
	assert.Greater(t, len(p2.Options), 0)
	assert.Greater(t, len(p1.Options), len(p2.Options))
}

func TestMultiSelectFilterCursor(t *testing.T) {
	p := newMultiSelectPrompt()
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

func TestMultiSelectFilterOutOptions(t *testing.T) {
	p := newMultiSelectPrompt()
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

	p.PressKey(&core.Key{Char: "#"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 0, p.CursorIndex)
}

func TestMultiSelectFilterOverSelectAll(t *testing.T) {
	p := newMultiSelectPrompt()
	p.Filter = true

	assert.Equal(t, 0, len(p.Value))
	assert.Equal(t, 3, len(p.Options))

	p.PressKey(&core.Key{Name: "a", Char: "a"})
	assert.Equal(t, 0, len(p.Value))
	assert.Equal(t, 2, len(p.Options))
	assert.Equal(t, "a", p.Search)
}

func TestMultiSelectInvalidOptionSelection(t *testing.T) {
	p := newMultiSelectPrompt()
	p.Filter = true

	p.PressKey(&core.Key{Name: "#", Char: "#"})
	assert.Equal(t, 0, len(p.Options))
	assert.Equal(t, 0, len(p.Value))

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, 0, len(p.Options))
	assert.Equal(t, 0, len(p.Value))
}
