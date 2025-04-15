package core

// Action represents an action that can be performed in the application.
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

// Custom messages for prompts
type SettingsMessages struct {
	// Custom message to display when a spinner is cancelled (default: "Canceled").
	CancelMessage string
	// Custom message to display when a spinner encounters an error (default: "Something went wrong").
	ErrorMessage string
}

// SettingsOptions defines user-configurable Settings for the application.
type SettingsOptions struct {
	// Aliases are custom key bindings for actions.
	// If a key binding already exists in the aliases map, it is not overwritten.
	Aliases map[KeyName]Action
	// Messages contains custom messages for the application.
	Messages SettingsMessages
}

var Settings = SettingsOptions{
	// aliases is a map that associates KeyName values with their corresponding Action.
	// It defines default key bindings for actions in the application.
	// Within any new alias coming from the user's land
	Aliases: map[KeyName]Action{
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
	},
	// Messages contains default messages for the application.
	Messages: SettingsMessages{
		CancelMessage: "Canceled",
		ErrorMessage:  "Something went wrong",
	},
}

// UpdateSettings updates the global SettingsOptions for the application.
func UpdateSettings(updates SettingsOptions) {
	for alias, action := range updates.Aliases {
		if _, exists := Settings.Aliases[alias]; !exists {
			Settings.Aliases[alias] = action
		}
	}

	if updates.Messages.CancelMessage != "" {
		Settings.Messages.CancelMessage = updates.Messages.CancelMessage
	}
	if updates.Messages.ErrorMessage != "" {
		Settings.Messages.ErrorMessage = updates.Messages.ErrorMessage
	}
}

// NewActionHandler creates a closure that handles key events and maps them to actions.
// It uses the global aliases map to determine the action for a given key and invokes the corresponding listener.
// If no listener is found for the action, the default listener is invoked.
//
// Parameters:
//   - listeners (map[Action]func()): A map of actions to their corresponding listener functions.
//   - defaultListener (func(key *Key)): A default listener function to invoke if no action-specific listener is found.
//
// Returns:
//   - actionHandler (func(key *Key)): A action handler that handles key events and invokes the appropriate listener.
func NewActionHandler(listeners map[Action]func(), defaultListener func(key *Key)) (actionHandler func(key *Key)) {
	return func(key *Key) {
		if action, actionExists := Settings.Aliases[key.Name]; actionExists {
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
