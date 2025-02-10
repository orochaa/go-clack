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

func runSelectKey() {
	prompts.SelectKey(prompts.SelectKeyParams[string]{
		Message: message,
		Options: []prompts.SelectKeyOption[string]{
			{Key: "f", Label: "Foo"},
			{Key: "b", Label: "Bar"},
			{Key: "Enter", Label: "Baz"},
		},
	})
}

func TestSelectKeyInitialState(t *testing.T) {
	go runSelectKey()
	time.Sleep(time.Millisecond)
	p := test.SelectKeyTestingPrompt.(*core.SelectKeyPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectKeyCancelState(t *testing.T) {
	go runSelectKey()
	time.Sleep(time.Millisecond)

	p := test.SelectKeyTestingPrompt.(*core.SelectKeyPrompt[string])
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectKeySubmitState(t *testing.T) {
	go runSelectKey()
	time.Sleep(time.Millisecond)

	p := test.SelectKeyTestingPrompt.(*core.SelectKeyPrompt[string])
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
