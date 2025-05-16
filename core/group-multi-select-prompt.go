package core

import (
	"context"
	"fmt"
	"os"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/core/validator"
)

type GroupMultiSelectOption[TValue comparable] struct {
	MultiSelectOption[TValue]
	IsGroup bool
	Options []*GroupMultiSelectOption[TValue]
}

func (o *GroupMultiSelectOption[TValue]) String() string {
	if o == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"{Label:%s, Value:%v(%T), Hint:%s, IsSelected:%t, IsGroup:%t, Options:%d}",
		o.Label,
		o.Value, o.Value,
		o.Hint,
		o.IsSelected,
		o.IsGroup,
		len(o.Options),
	)
}

type GroupMultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options        []*GroupMultiSelectOption[TValue]
	DisabledGroups bool
	Required       bool
}

type GroupMultiSelectPromptParams[TValue comparable] struct {
	Context        context.Context
	Input          *os.File
	Output         *os.File
	Options        map[string][]MultiSelectOption[TValue]
	InitialValue   []TValue
	DisabledGroups bool
	Required       bool
	Validate       func(value []TValue) error
	Render         func(p *GroupMultiSelectPrompt[TValue]) string
}

// NewGroupMultiSelectPrompt initializes and returns a new instance of GroupMultiSelectPrompt.
//
// The user can navigate between options using arrow keys.
// The user can select multiple options using space key.
// The prompt returns the values of the selected options.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Options (map[string][]MultiSelectOption[TValue]): A map of grouped options for the prompt (default: nil.
//   - InitialValue ([]TValue): The initial selected values (default: nil.
//   - DisabledGroups (bool): Whether groups are disabled for selection (default: false).
//   - Required (bool): Whether the prompt requires at least one selection (default: false).
//   - Validate (func(value []TValue) error): Custom validation function for the prompt (default: nil).
//   - Render (func(p *GroupMultiSelectPrompt[TValue]) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *GroupMultiSelectPrompt[TValue]: A new instance of GroupMultiSelectPrompt.
func NewGroupMultiSelectPrompt[TValue comparable](params GroupMultiSelectPromptParams[TValue]) *GroupMultiSelectPrompt[TValue] {
	v := validator.NewValidator("GroupMultiSelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	options := mapGroupMultiSelectOptions(params.Options)

	var p GroupMultiSelectPrompt[TValue]
	p = GroupMultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: mapGroupMultiSelectInitialValue(params.InitialValue, options),
			Validate:     WrapValidate(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]TValue](&p, params.Render),
		}),
		Options:        options,
		DisabledGroups: params.DisabledGroups,
		Required:       params.Required,
	}

	if p.DisabledGroups {
		p.CursorIndex = 1
	}

	actionHandler := NewActionHandler(map[Action]func(){
		UpAction:    func() { p.moveCursor(-1) },
		DownAction:  func() { p.moveCursor(1) },
		LeftAction:  func() { p.moveCursor(-1) },
		RightAction: func() { p.moveCursor(1) },
		HomeAction:  func() { p.CursorIndex = 0 },
		EndAction:   func() { p.CursorIndex = len(p.Options) - 1 },
		SpaceAction: p.toggleOption,
	}, nil)
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))
	})

	return &p
}

// IsGroupSelected checks if all options within a group are selected.
// If groups are disabled, it always returns false.
//
// Parameters:
//   - group (*GroupMultiSelectOption[TValue]): The group to check for selection.
//
// Returns:
//   - bool: True if all options in the group are selected, false otherwise.
func (p *GroupMultiSelectPrompt[TValue]) IsGroupSelected(group *GroupMultiSelectOption[TValue]) bool {
	if p.DisabledGroups {
		return false
	}
	for _, option := range group.Options {
		if !option.IsSelected {
			return false
		}

	}
	return true
}

// moveCursor moves the cursor up or down within the list of options.
// If groups are disabled, it skips over group headers.
//
// Parameters:
//   - direction (int): The direction to move the cursor (-1 for up, 1 for down).
func (p *GroupMultiSelectPrompt[TValue]) moveCursor(direction int) {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
	if p.DisabledGroups && p.Options[p.CursorIndex].IsGroup {
		p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
	}
}

// toggleOption toggles the selection state of the currently highlighted option or group.
// If a group is selected, it toggles all options within the group.
// If an option is selected, it updates the prompt's value accordingly.
func (p *GroupMultiSelectPrompt[TValue]) toggleOption() {
	option := p.Options[p.CursorIndex]

	if option.IsGroup && p.IsGroupSelected(option) {
		for _, option := range option.Options {
			option.IsSelected = false
		}
		p.Value = make([]TValue, 0, len(p.Value)-len(option.Options))
		for _, option := range p.Options {
			if option.IsSelected {
				p.Value = append(p.Value, option.Value)
			}
		}
		return
	}

	if option.IsGroup {
		for _, option := range option.Options {
			if !option.IsSelected {
				option.IsSelected = true
				p.Value = append(p.Value, option.Value)
			}
		}
		return
	}

	if option.IsSelected {
		option.IsSelected = false
		for i, v := range p.Value {
			if v == option.Value {
				p.Value = append(p.Value[:i], p.Value[i+1:]...)
				break
			}
		}
		return
	}

	option.IsSelected = true
	p.Value = append(p.Value, option.Value)
}

// mapGroupMultiSelectOptions converts a map of grouped options into a flat slice of GroupMultiSelectOption.
// It ensures that group headers and options are properly structured for rendering and selection.
//
// Parameters:
//   - groups (map[string][]MultiSelectOption[TValue]): A map of grouped options.
//
// Returns:
//   - []*GroupMultiSelectOption[TValue]: A flat slice of GroupMultiSelectOption.
func mapGroupMultiSelectOptions[TValue comparable](groups map[string][]MultiSelectOption[TValue]) []*GroupMultiSelectOption[TValue] {
	var options []*GroupMultiSelectOption[TValue]

	for groupName, groupOptions := range groups {
		group := &GroupMultiSelectOption[TValue]{
			MultiSelectOption: MultiSelectOption[TValue]{
				Label: groupName,
			},
			IsGroup: true,
			Options: make([]*GroupMultiSelectOption[TValue], len(groupOptions)),
		}
		options = append(options, group)
		for i, groupOption := range groupOptions {
			option := &GroupMultiSelectOption[TValue]{
				MultiSelectOption: MultiSelectOption[TValue]{
					Label:      groupOption.Label,
					Value:      groupOption.Value,
					IsSelected: groupOption.IsSelected,
				},
			}
			if value, ok := any(option.Value).(string); ok && value == "" {
				option.Value = any(option.Label).(TValue)
			}
			group.Options[i] = option
			options = append(options, option)
		}
	}

	return options
}

// mapGroupMultiSelectInitialValue maps the initial selected values to the corresponding options.
// If no initial values are provided, it uses the default selected options.
//
// Parameters:
//   - value ([]TValue): The initial selected values.
//   - options ([]*GroupMultiSelectOption[TValue]): The list of options to map the values to.
//
// Returns:
//   - []TValue: The mapped initial values.
func mapGroupMultiSelectInitialValue[TValue comparable](value []TValue, options []*GroupMultiSelectOption[TValue]) []TValue {
	if len(value) > 0 {
		return value
	}

	var initialValue []TValue
	for _, option := range options {
		if option.IsSelected {
			initialValue = append(initialValue, option.Value)
		}
	}
	return initialValue
}
