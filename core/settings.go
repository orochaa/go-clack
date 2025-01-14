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

// NewActionHandler creates a closure of a action handler based on the provided key and actions map
func NewActionHandler(listeners map[Action]func(), defaultListener func(key *Key)) (actionHandler func(key *Key)) {
	return func(key *Key) {
		if action, actionExists := aliases[key.Name]; actionExists {
			if listener, listenerExists := listeners[action]; listenerExists {
				if listener != nil {
					listener()
				}
				return
			}
		}
		if defaultListener != nil {
			defaultListener(key)
		}
	}
}
