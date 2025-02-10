package core_test

import (
	"errors"
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func TestWrapRender(t *testing.T) {
	prompt := core.NewPrompt(core.PromptParams[string]{
		InitialValue: "value",
	})
	render := func(p *core.Prompt[string]) string {
		return "frame: " + p.Value
	}
	renderWrapper := core.WrapRender[string](prompt, render)

	assert.Equal(t, "frame: value", renderWrapper(nil))
}

func TestWrapValidate(t *testing.T) {
	requiredErr := errors.New("required error")
	validateErr := errors.New("validate error")
	testStruct := struct{ a string }{a: "a"}

	type CustomType struct {
		Field string
	}

	testCases := []struct {
		description string
		required    bool
		validate    func(value any) error
		value       any
		expected    error
	}{
		{
			description: "validate without error",
			validate:    func(value any) error { return nil },
			expected:    nil,
		},
		{
			description: "validate with error",
			validate:    func(value any) error { return validateErr },
			expected:    validateErr,
		},
		{
			description: "required error",
			required:    true,
			expected:    requiredErr,
		},
		{
			description: "required validate without error",
			validate:    func(value any) error { return nil },
			required:    true,
			expected:    requiredErr,
		},
		{
			description: "required validate with error",
			validate:    func(value any) error { return validateErr },
			required:    true,
			expected:    validateErr,
		},
		{
			description: "required: true, value: false boolean",
			required:    true,
			value:       false,
			expected:    nil,
		},
		{
			description: "required: true, value: true boolean",
			required:    true,
			value:       true,
			expected:    nil},
		{
			description: "required: true, value: integer 0",
			required:    true,
			value:       0,
			expected:    requiredErr,
		},
		{
			description: "required: true, value: integer 1",
			required:    true,
			value:       1,
			expected:    nil,
		},
		{
			description: "required: true, value: empty string",
			required:    true,
			value:       "",
			expected:    requiredErr,
		},
		{
			description: "required: true, value: non-empty string",
			required:    true,
			value:       "a",
			expected:    nil,
		},
		{
			description: "required: true, value: empty slice",
			required:    true,
			value:       []any{},
			expected:    requiredErr,
		},
		{
			description: "required: true, value: non-empty slice",
			required:    true,
			value:       []any{"a"},
			expected:    nil,
		},
		{
			description: "required: true, value: nil value",
			required:    true,
			value:       nil,
			expected:    requiredErr,
		},
		{
			description: "required: true, value: empty struct",
			required:    true,
			value:       struct{}{},
			expected:    requiredErr,
		},
		{
			description: "required: true, value: non-empty struct",
			required:    true,
			value:       struct{ a string }{a: "a"},
			expected:    nil,
		},
		{
			description: "required: true, value: nil pointer",
			required:    true,
			value:       (*any)(nil),
			expected:    requiredErr,
		},
		{
			description: "required: true, value: valid pointer",
			required:    true,
			value:       &testStruct,
			expected:    nil,
		},
		{
			description: "required: true, value: empty array",
			required:    true,
			value:       [0]any{},
			expected:    requiredErr,
		},
		{
			description: "required: true, value: array with nil",
			required:    true,
			value:       [1]any{},
			expected:    requiredErr,
		},
		{
			description: "required: true, value: array with value",
			required:    true,
			value:       [1]any{1},
			expected:    nil,
		},
		{
			description: "required: true, value: empty map",
			required:    true,
			value:       map[string]string{},
			expected:    requiredErr,
		},
		{
			description: "required: true, value: non-empty map",
			required:    true,
			value:       map[string]string{"a": "b"},
			expected:    nil,
		},
		{
			description: "required: true, value: channel",
			required:    true,
			value:       make(chan any),
			expected:    nil,
		},
		{
			description: "required: true, value: empty custom type",
			required:    true,
			value:       CustomType{},
			expected:    requiredErr,
		},
		{
			description: "required: true, value: non-empty custom type",
			required:    true,
			value:       CustomType{Field: "a"},
			expected:    nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			validate := core.WrapValidate(tC.validate, &tC.required, requiredErr.Error())
			assert.Equal(t, tC.expected, validate(tC.value))
		})
	}
}
