package core

import (
	"os"
	"regexp"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type MultiSelectOption[TValue comparable] struct {
	Label      string
	Value      TValue
	IsSelected bool
}

type MultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	initialOptions []*MultiSelectOption[TValue]
	Options        []*MultiSelectOption[TValue]
	Search         string
	Filter         bool
	Required       bool
}

type MultiSelectPromptParams[TValue comparable] struct {
	Input        *os.File
	Output       *os.File
	InitialValue []TValue
	Options      []*MultiSelectOption[TValue]
	Filter       bool
	Required     bool
	Validate     func(value []TValue) error
	Render       func(p *MultiSelectPrompt[TValue]) string
}

func NewMultiSelectPrompt[TValue comparable](params MultiSelectPromptParams[TValue]) *MultiSelectPrompt[TValue] {
	v := validator.NewValidator("MultiSelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	for _, option := range params.Options {
		if value, ok := any(option.Value).(string); ok && value == "" {
			option.Value = any(option.Label).(TValue)
		}
	}

	var p MultiSelectPrompt[TValue]
	p = MultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: mapMultiSelectInitialValue(params.InitialValue, params.Options),
			Validate:     WrapValidate(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]TValue](&p, params.Render),
		}),
		initialOptions: params.Options,
		Options:        params.Options,
		Filter:         params.Filter,
		Required:       params.Required,
	}

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

func (p *MultiSelectPrompt[TValue]) handleKeyPress(key *Key) {
	moveCursor := func(direction int) {
		p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
	}

	HandleKeyAction(key, map[Action]func(){
		UpAction:    func() { moveCursor(-1) },
		DownAction:  func() { moveCursor(1) },
		LeftAction:  func() { moveCursor(-1) },
		RightAction: func() { moveCursor(1) },
		HomeAction:  func() { p.CursorIndex = 0 },
		EndAction:   func() { p.CursorIndex = len(p.Options) - 1 },
		SpaceAction: func() {
			if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
				option := p.Options[p.CursorIndex]
				if option.IsSelected {
					option.IsSelected = false
					value := []TValue{}
					for _, v := range p.Value {
						if v != option.Value {
							value = append(value, v)
						}
					}
					p.Value = value
					return
				}

				option.IsSelected = true
				p.Value = append(p.Value, option.Value)
			}
		},
		SubmitAction: nil,
		CancelAction: nil,
		DefaultAction: func() {
			if p.Filter {
				p.filterOptions(key)
			} else if key.Name == "a" {
				p.selectAll()
			}
		},
	})
}

func (p *MultiSelectPrompt[TValue]) selectAll() {
	if len(p.Value) == len(p.Options) {
		p.Value = []TValue{}
		for _, option := range p.Options {
			option.IsSelected = false
		}
		return
	}

	p.Value = make([]TValue, len(p.Options))
	for i, option := range p.Options {
		option.IsSelected = true
		p.Value[i] = option.Value
	}
}

func (p *MultiSelectPrompt[TValue]) filterOptions(key *Key) {
	var currentOption *MultiSelectOption[TValue]
	if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
		currentOption = p.Options[p.CursorIndex]
	}

	p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
	p.CursorIndex = 0

	if p.Search == "" {
		p.Options = p.initialOptions
		if currentOption == nil {
			return
		}

		for i, option := range p.Options {
			if option.Value == currentOption.Value {
				p.CursorIndex = i
			}
		}
		return
	}

	p.Options = []*MultiSelectOption[TValue]{}
	searchRegex, err := regexp.Compile("(?i)" + p.Search)
	if err != nil {
		return
	}
	for _, option := range p.initialOptions {
		if matched := searchRegex.MatchString(option.Label); matched {
			p.Options = append(p.Options, option)
			if currentOption != nil && option.Value == currentOption.Value {
				p.CursorIndex = len(p.Options) - 1
			}
		}
	}
}

func mapMultiSelectInitialValue[TValue comparable](value []TValue, options []*MultiSelectOption[TValue]) []TValue {
	if len(value) > 0 {
		for _, value := range value {
			for _, option := range options {
				if option.Value == value {
					option.IsSelected = true
				}
			}
		}
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
