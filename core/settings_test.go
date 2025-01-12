package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func TestTriggerActionWithKeyAlias(t *testing.T) {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	counter := 0
	core.HandleKeyAction(&core.Key{Name: "k"}, map[core.Action]func(){
		core.UpAction:   func() { counter++ },
		core.DownAction: t.FailNow,
	})

	assert.Equal(t, 1, counter)
}

func TestTriggerDefaultAction(t *testing.T) {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	counter := 0
	core.HandleKeyAction(&core.Key{Name: "l"}, map[core.Action]func(){
		core.UpAction:   t.FailNow,
		core.DownAction: t.FailNow,
		core.DefaultAction: func() {
			counter++
		},
	})

	assert.Equal(t, 1, counter)
}

func TestTriggerNoAction(t *testing.T) {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	core.HandleKeyAction(&core.Key{Name: "l"}, map[core.Action]func(){
		core.UpAction:   t.FailNow,
		core.DownAction: t.FailNow,
	})
}

func TestTriggerAliasActionOverDefaultAction(t *testing.T) {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	counter := 0
	core.HandleKeyAction(&core.Key{Name: "k"}, map[core.Action]func(){
		core.UpAction:      func() { counter++ },
		core.DownAction:    t.FailNow,
		core.DefaultAction: t.FailNow,
	})

	assert.Equal(t, 1, counter)
}

func TestTriggerIgnoredActionOverDefaultAction(t *testing.T) {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	core.HandleKeyAction(&core.Key{Name: "k"}, map[core.Action]func(){
		core.UpAction:      nil,
		core.DownAction:    t.FailNow,
		core.DefaultAction: t.FailNow,
	})
}

func TestTriggerActionWithInternalKeyAlias(t *testing.T) {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			core.UpKey: core.DownAction,
		},
	})

	core.HandleKeyAction(&core.Key{Name: core.UpKey}, map[core.Action]func(){
		core.DownAction: t.FailNow,
	})
}
