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

func runSelectPath() {
	prompts.SelectPath(prompts.SelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
	})
}

func TestSelectPathInitialState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)
	p := test.SelectPathTestingPrompt

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathWithOptionChildren(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)
	p := test.SelectPathTestingPrompt

	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathCancelState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathSubmitState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathEmptyFilter(t *testing.T) {
	go prompts.SelectPath(prompts.SelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
		Filter:     true,
	})
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathFilledFilter(t *testing.T) {
	go prompts.SelectPath(prompts.SelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
		Filter:     true,
	})
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Char: "f"})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
