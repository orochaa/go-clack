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

func runSelect() {
	prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Options: []*prompts.SelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
	})
}

func TestSelectInitialState(t *testing.T) {
	go runSelect()
	time.Sleep(time.Millisecond)
	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectWithHint(t *testing.T) {
	go prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Options: []*prompts.SelectOption[string]{
			{Label: "foo", Hint: "hint-foo"},
			{Label: "bar", Hint: "hint-bar"},
			{Label: "baz", Hint: "hint-baz"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectCancelState(t *testing.T) {
	go runSelect()
	time.Sleep(time.Millisecond)

	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectSubmitState(t *testing.T) {
	go runSelect()
	time.Sleep(time.Millisecond)

	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectWithLongList(t *testing.T) {
	go prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Options: []*prompts.SelectOption[string]{
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "d"},
			{Label: "e"},
			{Label: "f"},
			{Label: "g"},
			{Label: "h"},
			{Label: "i"},
			{Label: "j"},
			{Label: "k"},
			{Label: "l"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectEmptyFilter(t *testing.T) {
	go prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Filter:  true,
		Options: []*prompts.SelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
	})
	time.Sleep(time.Millisecond)

	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectFilledFilter(t *testing.T) {
	go prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Filter:  true,
		Options: []*prompts.SelectOption[string]{
			{Label: "foo"},
			{Label: "bar"},
			{Label: "baz"},
		},
	})
	time.Sleep(time.Millisecond)

	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])
	p.PressKey(&core.Key{Char: "b"})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
