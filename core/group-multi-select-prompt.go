package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type GroupMultiSelectOption[TValue comparable] struct {
	MultiSelectOption[TValue]
	IsGroup bool
	Options []*GroupMultiSelectOption[TValue]
}

type GroupMultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options        []*GroupMultiSelectOption[TValue]
	DisabledGroups bool
	Required       bool
}

type GroupMultiSelectPromptParams[TValue comparable] struct {
	Input          *os.File
	Output         *os.File
	Options        map[string][]MultiSelectOption[TValue]
	InitialValue   []TValue
	DisabledGroups bool
	Required       bool
	Validate       func(value []TValue) error
	Render         func(p *GroupMultiSelectPrompt[TValue]) string
}

func NewGroupMultiSelectPrompt[TValue comparable](params GroupMultiSelectPromptParams[TValue]) *GroupMultiSelectPrompt[TValue] {
	v := validator.NewValidator("GroupMultiSelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	options := mapGroupMultiSelectOptions(params.Options)

	var p GroupMultiSelectPrompt[TValue]
	p = GroupMultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
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

func (p *GroupMultiSelectPrompt[TValue]) moveCursor(direction int) {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
	if p.DisabledGroups && p.Options[p.CursorIndex].IsGroup {
		p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
	}
}

func (p *GroupMultiSelectPrompt[TValue]) toggleOption() {
	option := p.Options[p.CursorIndex]

	if option.IsGroup && p.IsGroupSelected(option) {
		for _, option := range option.Options {
			option.IsSelected = false
		}
		p.Value = []TValue{}
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
		p.Value = []TValue{}
		for _, _option := range p.Options {
			if _option.IsSelected {
				p.Value = append(p.Value, _option.Value)
			}
		}
		return
	}

	option.IsSelected = true
	p.Value = append(p.Value, option.Value)
}

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
