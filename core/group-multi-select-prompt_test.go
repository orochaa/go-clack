package core_test

import (
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newGroupMultiSelectPrompt() *core.GroupMultiSelectPrompt[string] {
	return core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Options: map[string][]core.MultiSelectOption[string]{
			"g1": {
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
			},
			"g2": {
				{Value: "x"},
				{Value: "y"},
				{Value: "z"},
			},
		},
		Render: func(p *core.GroupMultiSelectPrompt[string]) string {
			return ""
		},
	})
}

func TestChangeGroupMultiSelectCursor(t *testing.T) {
	p := newGroupMultiSelectPrompt()

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

func TestSelectGroupMultiSelectGroup(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	assert.Len(t, p.Value, 0)

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Len(t, p.Value, 3)

	p.CursorIndex = 4
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Len(t, p.Value, 6)

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Len(t, p.Value, 3)
}

func TestSelectGroupMultiSelectOption(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	assert.Len(t, p.Value, 0)

	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)

	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value, p.Options[2].Value}, p.Value)

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)
}

func TestGroupMultiSelectIsGroupSelected(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	group := p.Options[0]
	assert.Equal(t, false, p.IsGroupSelected(group))

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, true, p.IsGroupSelected(group))

	p.DisabledGroups = true
	assert.Equal(t, false, p.IsGroupSelected(group))
}

func TestLabelAsGroupMultiSelectValue(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Options: map[string][]core.MultiSelectOption[string]{
			"g1": {
				{Label: "a"},
				{Label: "b"},
				{Label: "c"},
			},
		},
		Render: func(p *core.GroupMultiSelectPrompt[string]) string {
			return ""
		},
	})

	assert.NotEmpty(t, p.Options[1].Value)

	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)
}

func TestGroupMultiSelectInitialValue(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		InitialValue: []string{"a", "b", "c"},
		Options: map[string][]core.MultiSelectOption[string]{
			"g1": {
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
			},
			"g2": {
				{Value: "x"},
				{Value: "y"},
				{Value: "z"},
			},
		},
		Render: func(p *core.GroupMultiSelectPrompt[string]) string {
			return ""
		},
	})

	assert.Len(t, p.Value, 3)
}

func TestGroupMultiSelectInitialValueAsIsSelected(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Options: map[string][]core.MultiSelectOption[string]{
			"g1": {
				{Value: "a"},
				{Value: "b", IsSelected: true},
				{Value: "c"},
			},
			"g2": {
				{Value: "x", IsSelected: true},
				{Value: "y", IsSelected: true},
				{Value: "z", IsSelected: true},
			},
		},
		Render: func(p *core.GroupMultiSelectPrompt[string]) string {
			return ""
		},
	})

	assert.Equal(t, []string{"b", "x", "y", "z"}, p.Value)
}

func TestGroupMultiSelectInitialValueOverIsSelected(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		InitialValue: []string{"a", "b", "c"},
		Options: map[string][]core.MultiSelectOption[string]{
			"g1": {
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
			},
			"g2": {
				{Value: "x", IsSelected: true},
				{Value: "y", IsSelected: true},
				{Value: "z"},
			},
		},
		Render: func(p *core.GroupMultiSelectPrompt[string]) string {
			return ""
		},
	})

	assert.Equal(t, []string{"a", "b", "c"}, p.Value)
}

func TestGroupMultiSelectRequiredValue(t *testing.T) {
	p := newGroupMultiSelectPrompt()
	p.Required = true

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}

func TestGroupMultiSelectDisabledGroups(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		DisabledGroups: true,
		Options: map[string][]core.MultiSelectOption[string]{
			"foo": {
				{Value: "a"},
				{Value: "b"},
			},
			"bar": {
				{Value: "x"},
				{Value: "y"},
			},
		},
		Render: func(p *core.GroupMultiSelectPrompt[string]) string { return "" },
	})

	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 4, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 5, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, 5, p.CursorIndex)
}
