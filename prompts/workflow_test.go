package prompts_test

import (
	"errors"
	"testing"

	"github.com/orochaa/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowStep(t *testing.T) {
	var r struct {
		Name string
	}

	prompts.Workflow(&r).
		Step("Name", func() (any, error) {
			return "test name", nil
		}).
		Run()

	assert.Equal(t, "test name", r.Name)
}

func TestWorkflowConditionalStep(t *testing.T) {
	var r struct {
		Project string
		Tools   []string
	}
	var project string
	w := prompts.Workflow(&r).
		Step("Project", func() (any, error) {
			return project, nil
		}).
		ConditionalStep("Tools",
			func() bool { return r.Project == "Node" },
			func() (any, error) { return []string{"Eslint"}, nil },
		)

	project = "Go"
	w.Run()
	assert.Equal(t, []string(nil), r.Tools)

	project = "Node"
	w.Run()
	assert.Equal(t, []string{"Eslint"}, r.Tools)
}

func TestWorkflowForkStep(t *testing.T) {
	type PackageManager struct {
		Command string
		Args    []string
	}
	var r struct {
		Install        bool
		PackageManager PackageManager
	}
	var install bool
	w := prompts.Workflow(&r).
		Step("Install", func() (any, error) {
			return install, nil
		}).
		ForkStep("PackageManager",
			func() bool { return r.Install },
			func() *prompts.WorkflowBuilder {
				return prompts.Workflow(&r.PackageManager).
					Step("Command", func() (any, error) {
						return "npm", nil
					}).
					Step("Args", func() (any, error) {
						return []string{"install"}, nil
					})
			},
		)

	install = true
	w.Run()
	assert.Equal(t, true, r.Install)
	assert.Equal(t, "npm", r.PackageManager.Command)
	assert.Equal(t, []string{"install"}, r.PackageManager.Args)

	install = false
	r.PackageManager = PackageManager{}
	w.Run()
	assert.Equal(t, false, r.Install)
	assert.Equal(t, "", r.PackageManager.Command)
	assert.Equal(t, []string(nil), r.PackageManager.Args)
}

func TestWorkflowLogStep(t *testing.T) {
	var r struct{}
	var calledTimes int

	err := prompts.Workflow(&r).
		LogStep("Name", func() {
			calledTimes++
		}).
		Run()

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, calledTimes)
}

func TestWorkflowStepError(t *testing.T) {
	var r struct {
		Name string
		Age  int
	}
	var calledTimes int

	err := prompts.Workflow(&r).
		Step("Name", func() (any, error) {
			calledTimes++
			return "", errors.New("test error")
		}).
		Step("Age", func() (any, error) {
			assert.FailNow(t, "Age step should not be called")
			return 22, nil
		}).
		Run()

	assert.Equal(t, "test error", err.Error())
	assert.Equal(t, 1, calledTimes)
}

func TestWorkflowOnCancel(t *testing.T) {
	var r struct {
		Name string
		Age  int
	}
	var calledTimes int

	err := prompts.Workflow(&r).
		Step("Name", func() (any, error) {
			calledTimes++
			return "", errors.New("test error")
		}).
		Step("Age", func() (any, error) {
			assert.FailNow(t, "Age step should not be called")
			return 22, nil
		}).
		OnCancel(func(step string, err error) {
			calledTimes++
			assert.Equal(t, "test error", err.Error())
			assert.Equal(t, "Name", step)
		}).
		Run()

	assert.Equal(t, "test error", err.Error())
	assert.Equal(t, 2, calledTimes)
}

func TestWorkflowCaseInsensitive(t *testing.T) {
	var r struct {
		Name string
		Age  int
	}

	prompts.Workflow(&r).
		Step("name", func() (any, error) {
			return "test name", nil
		}).
		Step("age", func() (any, error) {
			return 22, nil
		}).
		Run()

	assert.Equal(t, "test name", r.Name)
	assert.Equal(t, 22, r.Age)
}
