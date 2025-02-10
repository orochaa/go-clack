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

func runGroupMultiSelect() {
	prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		Message: message,
		Options: map[string][]prompts.MultiSelectOption[string]{
			"1": {
				{Label: "a"},
				{Label: "b"},
				{Label: "c"},
			},
			"2": {
				{Label: "x"},
				{Label: "y"},
				{Label: "z"},
			},
		},
	})
}

func TestGroupMultiSelectInitialState(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectCancelState(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectSubmitState(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectWithLongList(t *testing.T) {
	go prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		Message: message,
		Options: map[string][]prompts.MultiSelectOption[string]{
			"1": {
				{Label: "a"},
				{Label: "b"},
				{Label: "c"},
				{Label: "d"},
				{Label: "e"},
				{Label: "f"},
			},
			"2": {
				{Label: "u"},
				{Label: "v"},
				{Label: "w"},
				{Label: "x"},
				{Label: "y"},
				{Label: "z"},
			},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectMultiValue(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.CursorIndex = 5
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectDisabledGroups(t *testing.T) {
	go prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		DisabledGroups: true,
		Message:        message,
		Options: map[string][]prompts.MultiSelectOption[string]{
			"1": {
				{Label: "a"},
				{Label: "b"},
			},
			"2": {
				{Label: "x"},
				{Label: "y"},
			},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	p.PressKey(&core.Key{Name: core.SpaceKey})
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectSpacedGroups(t *testing.T) {
	go prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		SpacedGroups: true,
		Message:      message,
		Options: map[string][]prompts.MultiSelectOption[string]{
			"1": {
				{Label: "a"},
				{Label: "b"},
			},
			"2": {
				{Label: "x"},
				{Label: "y"},
			},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	p.PressKey(&core.Key{Name: core.SpaceKey})
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectDisabledAndSpacedGroups(t *testing.T) {
	go prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		DisabledGroups: true,
		SpacedGroups:   true,
		Message:        message,
		Options: map[string][]prompts.MultiSelectOption[string]{
			"1": {
				{Label: "a"},
				{Label: "b"},
			},
			"2": {
				{Label: "x"},
				{Label: "y"},
			},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	p.PressKey(&core.Key{Name: core.SpaceKey})
	cupaloy.SnapshotT(t, p.Frame)
}
