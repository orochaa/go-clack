package core_test

import (
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func TestTriggerActionWithKeyAlias(t *testing.T) {
	core.UpdateSettings(core.SettingsOptions{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})
	counter := 0

	actionHandler := core.NewActionHandler(map[core.Action]func(){
		core.UpAction:   func() { counter++ },
		core.DownAction: t.FailNow,
	}, nil)
	actionHandler(&core.Key{Name: "k"})

	assert.Equal(t, 1, counter)
}

func TestTriggerDefaultAction(t *testing.T) {
	core.UpdateSettings(core.SettingsOptions{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})
	counter := 0

	actionHandler := core.NewActionHandler(map[core.Action]func(){
		core.UpAction:   t.FailNow,
		core.DownAction: t.FailNow,
	}, func(key *core.Key) {
		counter++
	})
	actionHandler(&core.Key{Name: "l"})

	assert.Equal(t, 1, counter)
}

func TestTriggerNoAction(t *testing.T) {
	core.UpdateSettings(core.SettingsOptions{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	actionHandler := core.NewActionHandler(map[core.Action]func(){
		core.UpAction:   t.FailNow,
		core.DownAction: t.FailNow,
	}, nil)
	actionHandler(&core.Key{Name: "l"})
}

func TestTriggerAliasActionOverDefaultAction(t *testing.T) {
	core.UpdateSettings(core.SettingsOptions{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})
	counter := 0

	actionHandler := core.NewActionHandler(map[core.Action]func(){
		core.UpAction:   func() { counter++ },
		core.DownAction: t.FailNow,
	}, func(key *core.Key) { t.FailNow() })
	actionHandler(&core.Key{Name: "k"})

	assert.Equal(t, 1, counter)
}

func TestTriggerIgnoredActionOverDefaultAction(t *testing.T) {
	core.UpdateSettings(core.SettingsOptions{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
		},
	})

	actionHandler := core.NewActionHandler(map[core.Action]func(){
		core.UpAction:   nil,
		core.DownAction: t.FailNow,
	}, func(key *core.Key) { t.FailNow() })
	actionHandler(&core.Key{Name: "k"})
}

func TestTriggerActionWithInternalKeyAlias(t *testing.T) {
	core.UpdateSettings(core.SettingsOptions{
		Aliases: map[core.KeyName]core.Action{
			core.UpKey: core.DownAction,
		},
	})

	actionHandler := core.NewActionHandler(map[core.Action]func(){
		core.DownAction: t.FailNow,
	}, nil)
	actionHandler(&core.Key{Name: core.UpKey})

}
