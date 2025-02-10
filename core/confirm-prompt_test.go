package core_test

import (
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newConfirmPrompt() *core.ConfirmPrompt {
	return core.NewConfirmPrompt(core.ConfirmPromptParams{
		Render: func(p *core.ConfirmPrompt) string {
			return ""
		},
	})
}

func TestChangeConfirmValue(t *testing.T) {
	p := newConfirmPrompt()

	assert.Equal(t, false, p.Value)
	assert.Equal(t, 0, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, true, p.Value)
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, false, p.Value)
	assert.Equal(t, 0, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, true, p.Value)
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, false, p.Value)
	assert.Equal(t, 0, p.CursorIndex)
}
