package core_test

import (
	"path/filepath"
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPathPrompt() *core.SelectPathPrompt {
	return core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Render: func(p *core.SelectPathPrompt) string { return "" },
	})
}

func TestSelectPathChangeCursor(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.EndKey})
	assert.Equal(t, len(p.Options())-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.HomeKey})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.HomeKey})
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, len(p.Options())-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 1, p.CursorIndex)
}

func TestSelectPathChangeValue(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, p.CurrentLayer[0].Path, p.Value)

	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)

	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.CurrentLayer[2].Path, p.Value)

	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)
}

func TestSelectPathEnterDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	p.CurrentOption.IsDir = true
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, 2, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption.Parent)
}

func TestSelectPathEnterNonDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	p.CurrentOption.IsDir = false
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
	assert.Equal(t, false, p.CurrentOption.IsOpen)
}

func TestSelectPathExitDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	pastOption := p.CurrentOption
	p.CurrentOption.Open()
	p.CurrentLayer = p.CurrentOption.Children
	p.CurrentOption = p.CurrentOption.Children[0]
	assert.Equal(t, 2, p.CurrentOption.Depth)

	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
	assert.Equal(t, false, p.CurrentOption.IsOpen)
}

func TestSelectPathExitEmptyDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	p.CurrentOption.IsDir = true
	p.CurrentOption.IsOpen = true
	p.CurrentOption.Children = []*core.PathNode{}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.LeftKey})

	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
	assert.Equal(t, false, p.CurrentOption.IsOpen)
}

func TestSelectPathExitRootDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, p.Root, p.CurrentOption)
	assert.Equal(t, true, p.CurrentOption.IsRoot())
	assert.Equal(t, true, p.CurrentOption.IsOpen)

	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, p.Root, p.CurrentOption)
	assert.Equal(t, true, p.CurrentOption.IsRoot())
	assert.Equal(t, true, p.CurrentOption.IsOpen)
	assert.Equal(t, filepath.Dir(pastOption.Path), p.CurrentOption.Path)
}

func TestSelectPathFilter(t *testing.T) {
	p1 := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Render: func(p *core.SelectPathPrompt) string { return "" },
	})
	p2 := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Filter: true,
		Render: func(p *core.SelectPathPrompt) string { return "" },
	})

	p1.PressKey(&core.Key{Char: "d"})
	p2.PressKey(&core.Key{Char: "d"})

	assert.Greater(t, len(p1.Options()), 0)
	assert.Greater(t, len(p2.Options()), 0)
	assert.Greater(t, len(p1.Options()), len(p2.Options()))
}

func TestSelectPathFilterNavigate(t *testing.T) {
	p := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Filter: true,
		Render: func(p *core.SelectPathPrompt) string { return "" },
	})

	p.PressKey(&core.Key{Char: "f"})
	assert.Equal(t, 3, p.CurrentOption.Index)

	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 4, p.CurrentOption.Index)
}

func TestSelectPathFilterRecover(t *testing.T) {
	p := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Filter:     true,
		FileSystem: MockFileSystem{},
		Render:     func(p *core.SelectPathPrompt) string { return "" },
	})

	assert.Equal(t, 3, len(p.Options()))

	p.PressKey(&core.Key{Char: "z"})
	assert.Equal(t, 1, len(p.Options()))

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 3, len(p.Options()))
}
