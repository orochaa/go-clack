package prompts_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/stretchr/testify/assert"
)

func TestTasksStart(t *testing.T) {
	startTimes := 0
	task := func(message func(msg string)) (string, error) {
		startTimes++
		time.Sleep(time.Millisecond)
		return "", nil
	}
	w := &MockWriter{}

	prompts.Tasks(
		[]prompts.Task{
			{Title: "Foo", Task: task},
			{Title: "Bar", Task: task},
			{Title: "Baz", Task: task},
		}, prompts.SpinnerOptions{
			Output:        w,
			FrameInterval: time.Millisecond,
		})
	time.Sleep(5 * time.Millisecond)

	expectedList := []string{
		"◒ Foo",
		"◒ Bar",
		"◒ Baz",
	}
	for _, expected := range expectedList {
		assert.Contains(t, w.Data, expected)
	}
}

func TestTasksSubmit(t *testing.T) {
	startTimes := 0
	task := func(message func(msg string)) (string, error) {
		startTimes++
		time.Sleep(time.Millisecond)
		return "", nil
	}
	w := &MockWriter{}

	prompts.Tasks(
		[]prompts.Task{
			{Title: "Foo", Task: task},
			{Title: "Bar", Task: task},
			{Title: "Baz", Task: task},
		}, prompts.SpinnerOptions{
			Output:        w,
			FrameInterval: time.Millisecond,
		})
	time.Sleep(5 * time.Millisecond)

	expectedList := []string{
		symbols.STEP_SUBMIT + " Foo\n",
		symbols.STEP_SUBMIT + " Bar\n",
		symbols.STEP_SUBMIT + " Baz\n",
	}
	for _, expected := range expectedList {
		assert.Contains(t, w.Data, expected)
	}
}

func TestTasksUpdateMessage(t *testing.T) {
	task := func(message func(msg string)) (string, error) {
		message("Bar")
		time.Sleep(time.Millisecond)
		return "", nil
	}
	w := &MockWriter{}

	prompts.Tasks([]prompts.Task{
		{Title: "Foo", Task: task},
	}, prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})
	time.Sleep(2 * time.Millisecond)

	assert.Contains(t, w.Data, "◒ Bar")
}

func TestTasksWithDisabledTask(t *testing.T) {
	counter := 0
	task := func(message func(msg string)) (string, error) {
		counter++
		return "", nil
	}
	w := &MockWriter{}

	prompts.Tasks([]prompts.Task{
		{Title: "Foo", Task: task, Disabled: true},
	}, prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})
	time.Sleep(2 * time.Millisecond)

	assert.Equal(t, 0, counter)
}

func TestTasksTaskWithError(t *testing.T) {
	task := func(message func(msg string)) (string, error) {
		return "", errors.New("task error")
	}
	w := &MockWriter{}

	prompts.Tasks([]prompts.Task{
		{Title: "Foo", Task: task},
	}, prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})
	time.Sleep(2 * time.Millisecond)

	assert.Contains(t, w.Data, fmt.Sprintf("%s task error\n", symbols.STEP_CANCEL))
}
