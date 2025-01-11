package core_test

import (
	"path/filepath"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPathPrompt() *core.MultiSelectPathPrompt {
	return core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		Render: func(p *core.MultiSelectPathPrompt) string { return "" },
	})
}

func TestMultiSelectPathChangeCursor(t *testing.T) {
	p := newMultiSelectPathPrompt()

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

func TestMultiSelectPathChangeValue(t *testing.T) {
	p := newMultiSelectPathPrompt()
	options := p.Options()

	assert.Equal(t, []string(nil), p.Value)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{options[1].Path}, p.Value)
	assert.Equal(t, true, options[1].IsSelected)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{}, p.Value)
	assert.Equal(t, false, options[1].IsSelected)

	p.CurrentOption = options[1]
	for i := 0; i < 3; i++ {
		p.PressKey(&core.Key{Name: core.SpaceKey})
		p.PressKey(&core.Key{Name: core.DownKey})
	}
	p.CurrentOption = options[2]
	p.PressKey(&core.Key{Name: core.SpaceKey})
	expected := append([]string{options[1].Path}, options[3].Path)
	assert.Equal(t, expected, p.Value)
}

func TestMultiSelectPathEnterDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

	p.CurrentOption.IsDir = true
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, 2, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption.Parent)
}

func TestMultiSelectPathEnterNonDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

	p.CurrentOption.IsDir = false
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
	assert.Equal(t, false, p.CurrentOption.IsOpen)
}

func TestMultiSelectPathExitDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

	pastOption := p.CurrentOption
	p.CurrentOption.Open()
	p.CurrentOption = p.CurrentOption.Children[0]
	assert.Equal(t, 2, p.CurrentOption.Depth)

	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
	assert.Equal(t, false, p.CurrentOption.IsOpen)
}

func TestMultiSelectPathExitEmptyDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

	p.CurrentOption.IsDir = true
	p.CurrentOption.IsOpen = true
	p.CurrentOption.Children = []*core.PathNode{}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.LeftKey})

	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
	assert.Equal(t, false, p.CurrentOption.IsOpen)
}

func TestMultiSelectPathExitRootDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

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

func TestMultiSelectPathExitToSelectedDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()
	p.Value = []string{filepath.Dir(p.Root.Path)}

	p.PressKey(&core.Key{Name: core.LeftKey})
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.True(t, p.CurrentOption.IsSelected)
}

func TestMultiSelectPathRequiredValue(t *testing.T) {
	p := newMultiSelectPathPrompt()
	p.Required = true

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}

func TestMultiSelectPathFilter(t *testing.T) {
	p1 := core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		Render: func(p *core.MultiSelectPathPrompt) string { return "" },
	})
	p2 := core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		Filter:     true,
		FileSystem: MockFileSystem{},
		Render:     func(p *core.MultiSelectPathPrompt) string { return "" },
	})

	p1.PressKey(&core.Key{Char: "d"})
	p2.PressKey(&core.Key{Char: "d"})

	assert.Greater(t, len(p1.Options()), 0)
	assert.Greater(t, len(p2.Options()), 0)
	assert.Greater(t, len(p1.Options()), len(p2.Options()))
}

func TestMultiSelectPathFilterNavigate(t *testing.T) {
	p := core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		Filter: true,
		Render: func(p *core.MultiSelectPathPrompt) string { return "" },
	})

	p.PressKey(&core.Key{Char: "f"})
	assert.Equal(t, 3, p.CurrentOption.Index)

	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 4, p.CurrentOption.Index)
}

func TestMultiSelectPathFilterRecover(t *testing.T) {
	p := core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		Filter:     true,
		FileSystem: MockFileSystem{},
		Render:     func(p *core.MultiSelectPathPrompt) string { return "" },
	})

	assert.Equal(t, 3, len(p.Options()))

	p.PressKey(&core.Key{Char: "z"})
	assert.Equal(t, 1, len(p.Options()))

	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, 3, len(p.Options()))
}

func TestMultiSelectPathSortOnFinalize(t *testing.T) {
	p := newMultiSelectPathPrompt()

	p.Value = []string{"b", "a", "1"}
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, []string{"1", "a", "b"}, p.Value)
}
