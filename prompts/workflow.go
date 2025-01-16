package prompts

import (
	"fmt"
	"reflect"
	"strings"
)

type WorkflowStep struct {
	Name      string
	Prompt    func() (any, error)
	Condition func() bool
	SetResult bool
}

type WorkflowBuilder struct {
	steps    []*WorkflowStep
	onCancel func(step string, err error)
	result   any
}

// Step adds a new step to the workflow with a prompt function to gather the step's result.
func (w *WorkflowBuilder) Step(name string, prompt func() (any, error)) *WorkflowBuilder {
	w.steps = append(w.steps, &WorkflowStep{
		Name: name, Prompt: prompt,
		SetResult: true,
	})
	return w
}

// ConditionalStep adds a new step that runs only if the provided condition is true.
func (w *WorkflowBuilder) ConditionalStep(name string, condition func() bool, prompt func() (any, error)) *WorkflowBuilder {
	w.steps = append(w.steps, &WorkflowStep{
		Name:      name,
		Condition: condition,
		Prompt:    prompt,
		SetResult: true,
	})
	return w
}

// ForkStep adds a new step that runs a sub-workflow if the provided condition is true.
func (w *WorkflowBuilder) ForkStep(name string, condition func() bool, subWorkflow func() *WorkflowBuilder) *WorkflowBuilder {
	w.steps = append(w.steps, &WorkflowStep{
		Name:      name,
		Condition: condition,
		Prompt: func() (any, error) {
			err := subWorkflow().Run()
			return nil, err
		},
		SetResult: false,
	})
	return w
}

// LogStep adds a step that performs an action without affecting the workflow's result.
// Typically used for logging or other side effects.
func (w *WorkflowBuilder) LogStep(name string, action func()) *WorkflowBuilder {
	w.steps = append(w.steps, &WorkflowStep{
		Name: name,
		Prompt: func() (any, error) {
			action()
			return nil, nil
		},
		SetResult: false,
	})
	return w
}

// CustomStep adds a custom step to the workflow. The custom step can define its own behavior.
func (w *WorkflowBuilder) CustomStep(step *WorkflowStep) *WorkflowBuilder {
	w.steps = append(w.steps, step)
	return w
}

// OnCancel sets a callback function to be called if any step encounters an error and the workflow is onCanceled.
func (w *WorkflowBuilder) OnCancel(onCancel func(step string, err error)) *WorkflowBuilder {
	w.onCancel = onCancel
	return w
}

// Run executes all the steps in the workflow in sequence.
// If a step's condition is not met, it is skipped.
// If a step encounters an error, the onCancel callback is called and the error is returned.
func (w *WorkflowBuilder) Run() error {
	v := reflect.ValueOf(w.result).Elem()

	for _, step := range w.steps {
		if step.Condition != nil && !step.Condition() {
			continue
		}

		stepResult, stepErr := step.Prompt()
		if stepErr != nil {
			if w.onCancel != nil {
				w.onCancel(step.Name, stepErr)
			}
			return stepErr
		}
		if !step.SetResult {
			continue
		}

		field := v.FieldByName(capitalize(step.Name))
		if !field.IsValid() {
			panic(fmt.Sprintf("workflow error: invalid field `%s`", step.Name))
		}
		if !field.CanSet() {
			panic(fmt.Sprintf("workflow error: cannot set field `%s`", step.Name))
		}

		stepResultVal := reflect.ValueOf(stepResult)
		if !stepResultVal.Type().AssignableTo(field.Type()) {
			panic(fmt.Sprintf("workflow error: type `%s` is not assignable to field `%s`", stepResultVal.Type(), step.Name))
		}

		field.Set(stepResultVal)
	}
	return nil
}

// Workflow returns a new WorkflowBuilder instance, allowing you to define a workflow that populates the provided value.
// The value should be a struct, and each field will be populated by a corresponding step in the workflow.
// If the field does not correspond to a struct's field, it panics
func Workflow(v any) *WorkflowBuilder {
	return &WorkflowBuilder{result: v}
}

func capitalize(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(string(str[0])) + (str[1:])
}
