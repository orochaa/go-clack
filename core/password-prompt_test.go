package core_test

import (
	"fmt"
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPasswordPrompt() *core.PasswordPrompt {
	return core.NewPasswordPrompt(core.PasswordPromptParams{
		Render: func(p *core.PasswordPrompt) string { return "" },
	})
}

func TestPasswordPromptInitialValue(t *testing.T) {
	p := core.NewPasswordPrompt(core.PasswordPromptParams{
		InitialValue: "foo",
		Render:       func(p *core.PasswordPrompt) string { return "" },
	})

	assert.Equal(t, "foo", p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestChangePasswordValue(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, "", p.Value)
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)
	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, "a", p.Value)
}

func TestPasswordMask(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, "", p.ValueWithMask())
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "*", p.ValueWithMask())
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "**", p.ValueWithMask())
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, "**", p.ValueWithMask())
	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, "*", p.ValueWithMask())
}

func TestPasswordMaskWithCursor(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, "█", p.ValueWithMaskAndCursor())
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "*█", p.ValueWithMaskAndCursor())
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "**█", p.ValueWithMaskAndCursor())
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, "**", p.ValueWithMaskAndCursor())
	p.PressKey(&core.Key{Name: core.BackspaceKey})
	assert.Equal(t, "*", p.ValueWithMaskAndCursor())
}

func TestValidatePassword(t *testing.T) {
	p := core.NewPasswordPrompt(core.PasswordPromptParams{
		InitialValue: "123",
		Validate: func(value string) error {
			return fmt.Errorf("invalid password: %s", value)
		},
		Render: func(p *core.PasswordPrompt) string { return "" },
	})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, "invalid password: 123", p.Error)
}

func TestPasswordRequiredValue(t *testing.T) {
	p := newPasswordPrompt()
	p.Required = true

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}
