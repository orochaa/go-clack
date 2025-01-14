package core

import (
	"os"
	"regexp"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
}

type SelectPrompt[TValue comparable] struct {
	Prompt[TValue]
	initialOptions []*SelectOption[TValue]
	Options        []*SelectOption[TValue]
	Search         string
	Filter         bool
	Required       bool
}

type SelectPromptParams[TValue comparable] struct {
	Input        *os.File
	Output       *os.File
	InitialValue TValue
	Options      []*SelectOption[TValue]
	Filter       bool
	Required     bool
	Render       func(p *SelectPrompt[TValue]) string
}

func NewSelectPrompt[TValue comparable](params SelectPromptParams[TValue]) *SelectPrompt[TValue] {
	v := validator.NewValidator("SelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	startIndex := 0
	for i, option := range params.Options {
		if value, ok := any(option.Value).(string); ok && value == "" {
			option.Value = any(option.Label).(TValue)
		}
		if option.Value == params.InitialValue {
			startIndex = i
		}
	}

	var p SelectPrompt[TValue]
	p = SelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.Options[startIndex].Value,
			CursorIndex:  startIndex,
			Validate:     WrapValidate[TValue](nil, &p.Required, "Please select an option."),
			Render:       WrapRender[TValue](&p, params.Render),
		}),
		initialOptions: params.Options,
		Options:        params.Options,
		Filter:         params.Filter,
		Required:       params.Required,
	}

	actionHandler := NewActionHandler(map[Action]func(){
		UpAction:    func() { p.moveCursor(-1) },
		DownAction:  func() { p.moveCursor(1) },
		LeftAction:  func() { p.moveCursor(-1) },
		RightAction: func() { p.moveCursor(1) },
		HomeAction:  func() { p.CursorIndex = 0 },
		EndAction:   func() { p.CursorIndex = len(p.Options) - 1 },
	}, p.filterOptions)
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))

		if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
			p.Value = p.Options[p.CursorIndex].Value
		} else {
			p.Value = *new(TValue)
		}
	})

	return &p
}

func (p *SelectPrompt[TValue]) moveCursor(direction int) {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
}

func (p *SelectPrompt[TValue]) filterOptions(key *Key) {
	if !p.Filter {
		return
	}

	p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
	p.CursorIndex = 0

	if p.Search == "" {
		p.Options = p.initialOptions
		for i, option := range p.Options {
			if option.Value == p.Value {
				p.CursorIndex = i
				break
			}
		}
		return
	}

	p.Options = []*SelectOption[TValue]{}
	searchRegex, err := regexp.Compile("(?i)" + p.Search)
	if err != nil {
		return
	}

	for _, option := range p.initialOptions {
		if matched := searchRegex.MatchString(option.Label); matched {
			p.Options = append(p.Options, option)
			if option.Value == p.Value {
				p.CursorIndex = len(p.Options) - 1
			}
		}
	}
}
