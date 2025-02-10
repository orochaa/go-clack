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

func runMultiSelectPath() {
	prompts.MultiSelectPath(prompts.MultiSelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
	})
}

func TestMultiSelectPathInitialState(t *testing.T) {
	go runMultiSelectPath()
	time.Sleep(time.Millisecond)
	p := test.MultiSelectPathTestingPrompt

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectPathWithOptionChildren(t *testing.T) {
	go runMultiSelectPath()
	time.Sleep(time.Millisecond)
	p := test.MultiSelectPathTestingPrompt

	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectPathWithSelectedOptions(t *testing.T) {
	go runMultiSelectPath()
	time.Sleep(time.Millisecond)
	p := test.MultiSelectPathTestingPrompt

	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectPathCancelState(t *testing.T) {
	go runMultiSelectPath()
	time.Sleep(time.Millisecond)

	p := test.MultiSelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectPathSubmitState(t *testing.T) {
	go runMultiSelectPath()
	time.Sleep(time.Millisecond)

	p := test.MultiSelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectPathEmptyFilter(t *testing.T) {
	go prompts.MultiSelectPath(prompts.MultiSelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
		Filter:     true,
	})
	time.Sleep(time.Millisecond)

	p := test.MultiSelectPathTestingPrompt

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectPathFilledFilter(t *testing.T) {
	go prompts.MultiSelectPath(prompts.MultiSelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
		Filter:     true,
	})
	time.Sleep(time.Millisecond)

	p := test.MultiSelectPathTestingPrompt
	p.PressKey(&core.Key{Char: "f"})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
