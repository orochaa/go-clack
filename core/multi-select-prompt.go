package core

import (
	"context"
	"os"
	"regexp"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/core/validator"
)

type MultiSelectOption[TValue comparable] struct {
	Label      string
	Value      TValue
	Hint       string
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
	Context      context.Context
	Input        *os.File
	Output       *os.File
	Options      []*MultiSelectOption[TValue]
	InitialValue []TValue
	Filter       bool
	Required     bool
	Validate     func(value []TValue) error
	Render       func(p *MultiSelectPrompt[TValue]) string
}

// NewMultiSelectPrompt initializes and returns a new instance of MultiSelectPrompt.
//
// The user can navigate through options using arrow keys.
// The user can select multiple options using space key.
// The prompt returns the value of the selected options.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - Options ([]*MultiSelectOption[TValue]): A list of options for the prompt (default: nil.
//   - InitialValue ([]TValue): The initial selected values (default: nil.
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - Required (bool): Whether the prompt requires at least one selection (default: false).
//   - Validate (func(value []TValue) error): Custom validation function for the prompt (default: nil).
//   - Render (func(p *MultiSelectPrompt[TValue]) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *MultiSelectPrompt[TValue]: A new instance of MultiSelectPrompt.
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
			Context:      params.Context,
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

	actionHandler := NewActionHandler(map[Action]func(){
		UpAction:    func() { p.moveCursor(-1) },
		DownAction:  func() { p.moveCursor(1) },
		LeftAction:  func() { p.moveCursor(-1) },
		RightAction: func() { p.moveCursor(1) },
		HomeAction:  func() { p.CursorIndex = 0 },
		EndAction:   func() { p.CursorIndex = len(p.Options) - 1 },
		SpaceAction: p.toggleOption,
	}, func(key *Key) {
		if p.Filter {
			p.filterOptions(key)
		} else if key.Name == "a" {
			p.toggleAllOptions()
		}
	})
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))
	})

	return &p
}

// moveCursor moves the cursor up or down within the list of options.
// It ensures the cursor stays within the bounds of the available options.
//
// Parameters:
//   - direction (int): The direction to move the cursor (-1 for up, 1 for down).
func (p *MultiSelectPrompt[TValue]) moveCursor(direction int) {
	p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+direction, len(p.Options))
}

// toggleOption toggles the selection state of the currently highlighted option.
// If the option is selected, it is deselected, and vice versa.
// The prompt's value is updated accordingly.
func (p *MultiSelectPrompt[TValue]) toggleOption() {
	if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
		option := p.Options[p.CursorIndex]

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
}

// toggleAllOptions toggles the selection state of all options.
// If all options are selected, they are deselected, and vice versa.
// The prompt's value is updated accordingly.
func (p *MultiSelectPrompt[TValue]) toggleAllOptions() {
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

// filterOptions updates the search term based on the provided key input and filters the available options.
// If the search term is empty, it resets the options to the initial list.
//
// Parameters:
//   - key (*Key): The key event that triggered the filtering.
func (p *MultiSelectPrompt[TValue]) filterOptions(key *Key) {
	if !p.Filter {
		return
	}

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
				break
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

// mapMultiSelectInitialValue maps the initial selected values to the corresponding options.
// If no initial values are provided, it uses the default selected options.
//
// Parameters:
//   - value ([]TValue): The initial selected values.
//   - options ([]*MultiSelectOption[TValue]): The list of options to map the values to.
//
// Returns:
//   - []TValue: The mapped initial values.
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
