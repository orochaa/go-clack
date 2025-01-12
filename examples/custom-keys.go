package main

import (
	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
)

func CustomKeys() {
	core.UpdateSettings(core.Settings{
		Aliases: map[core.KeyName]core.Action{
			"k": core.UpAction,
			"j": core.DownAction,
			"h": core.LeftAction,
			"l": core.RightAction,
		},
	})

	prompts.SelectPath(prompts.SelectPathParams{
		Message: "Try Vim keys to move: (k=up,j=down,h=left,l=right)",
	})
}
