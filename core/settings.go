package core

type Action int

const (
	UpAction Action = iota
	DownAction
	LeftAction
	RightAction
	HomeAction
	EndAction
	SpaceAction
	SubmitAction
	CancelAction
	DefaultAction
)

var aliases = map[KeyName]Action{
	UpKey:     UpAction,
	DownKey:   DownAction,
	LeftKey:   LeftAction,
	RightKey:  RightAction,
	HomeKey:   HomeAction,
	EndKey:    EndAction,
	SpaceKey:  SpaceAction,
	EnterKey:  SubmitAction,
	CancelKey: CancelAction,
	EscapeKey: CancelAction,
}

// Settings defines user-configurable settings
type Settings struct {
	// Custom global aliases for the default actions
	Aliases map[KeyName]Action
}

// UpdateSettings updates the global settings based on the provided Settings
func UpdateSettings(settings Settings) {
	for alias, action := range settings.Aliases {
		if _, exists := aliases[alias]; !exists {
			aliases[alias] = action
		}
	}
}

// HandleKeyAction handles a key action based on the provided key and actions map
func HandleKeyAction(key *Key, actions map[Action]func()) {
	if action, actionExists := aliases[key.Name]; actionExists {
		if listener, listenerExists := actions[action]; listenerExists {
			if listener != nil {
				listener()
			}
			return
		}
	}
	if defaultListener := actions[DefaultAction]; defaultListener != nil {
		defaultListener()
	}
}
