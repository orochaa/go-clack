package prompts_test

import (
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/prompts"
	"github.com/orochaa/go-clack/prompts/test"
	"github.com/stretchr/testify/assert"
)

func runMultiSelect() {
	prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "foo", IsSelected: true},
			{Label: "bar", IsSelected: true},
			{Label: "baz"},
		},
	})
}

func TestMultiSelectInitialState(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectWithHint(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "foo", Hint: "hint-foo"},
			{Label: "bar", Hint: "hint-bar"},
			{Label: "baz", Hint: "hint-baz"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectWithSelectedHint(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "foo", Hint: "hint-foo"},
			{Label: "bar", Hint: "hint-bar"},
			{Label: "baz", Hint: "hint-baz"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	p.PressKey(&core.Key{Name: core.SpaceKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectCancelState(t *testing.T) {
	go runMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectSubmitState(t *testing.T) {
	go runMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectWithLongList(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectMultiValue(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "a", IsSelected: true},
			{Label: "b", IsSelected: true},
			{Label: "c"},
			{Label: "a", IsSelected: true},
			{Label: "b", IsSelected: true},
			{Label: "c"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: core.DownKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectEmptyFilter(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Filter:  true,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
	})
	time.Sleep(time.Millisecond)

	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectFilledFilter(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Filter:  true,
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
	})
	time.Sleep(time.Millisecond)

	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.PressKey(&core.Key{Char: "b"})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
