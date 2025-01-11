package core_test

import (
	"fmt"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPathPrompt() *core.PathPrompt {
	return core.NewPathPrompt(core.PathPromptParams{
		Render: func(p *core.PathPrompt) string { return "" },
	})
}

func TestPathDefaultValue(t *testing.T) {
	p := newPathPrompt()

	cwd, _ := p.FileSystem.Getwd()

	assert.Equal(t, cwd, p.Value)
	assert.Equal(t, len(cwd), p.CursorIndex)
}

func TestPathChangeValue(t *testing.T) {
	p := newPathPrompt()
	p.Prompt.Value = "/go-clack"
	p.Prompt.CursorIndex = len("/go-clack")
	p.Value = "/go-clack"

	assert.Equal(t, "/go-clack", p.Value)
	p.PressKey(&core.Key{Char: "/"})
	assert.Equal(t, "/go-clack/", p.Value)
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "/go-clack/a", p.Value)
	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, "/go-clack/", p.Value)
}

func TestPathChangeHint(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, "/", p.Hint)
	p.PressKey(&core.Key{Char: "/"})
	assert.Equal(t, "README.md", p.Hint)
	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, "/", p.Hint)
}

func TestPathTabComplete(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, 0, len(p.HintOptions))
	assert.Equal(t, -1, p.HintIndex)
	assert.Equal(t, "/", p.Hint)

	p.PressKey(&core.Key{Name: core.TabKey})
	assert.Equal(t, 0, len(p.HintOptions))
	assert.Equal(t, -1, p.HintIndex)
	assert.NotEqual(t, "", p.Hint)
	assert.NotEqual(t, "/", p.Hint)

	p.PressKey(&core.Key{Name: core.TabKey})
	assert.GreaterOrEqual(t, len(p.HintOptions), 1)
	assert.Equal(t, 0, p.HintIndex)

	p.PressKey(&core.Key{Name: core.TabKey})
	assert.GreaterOrEqual(t, len(p.HintOptions), 1)
	assert.Equal(t, 1, p.HintIndex)

	p.PressKey(&core.Key{Name: core.TabKey})
	assert.GreaterOrEqual(t, len(p.HintOptions), 1)
	assert.Equal(t, 2, p.HintIndex)
}

func TestPathComplete(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, "/", p.Hint)

	expected := p.Value + p.Hint
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, expected, p.Value)
	assert.NotEqual(t, "", p.Hint)
	assert.NotEqual(t, "/", p.Hint)

	expected = p.Value + p.Hint
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, expected, p.Value)
}

func TestPathValueWithHint(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, p.Value+p.Hint, p.ValueWithCursor())

	p.Hint = ""
	p.CursorIndex = len(p.Value)
	assert.Equal(t, p.Value+"█", p.ValueWithCursor())

	p.Hint = ""
	p.CursorIndex = len(p.Value) - 1
	assert.Equal(t, p.Value, p.ValueWithCursor())

	p.Hint = "/"
	p.CursorIndex = len(p.Value) - 1
	assert.Equal(t, p.Value+p.Hint, p.ValueWithCursor())
}

func TestPathValidateValue(t *testing.T) {
	p := core.NewPathPrompt(core.PathPromptParams{
		InitialValue: "/folder",
		Validate: func(value string) error {
			return fmt.Errorf("invalid path: %s", value)
		},
		Render: func(p *core.PathPrompt) string { return "" },
	})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, "invalid path: /folder", p.Error)
}

func TestPathRequiredValue(t *testing.T) {
	p := newPathPrompt()
	p.Required = true

	p.Value = ""
	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}
