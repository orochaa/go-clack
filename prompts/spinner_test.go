package prompts_test

import (
	"os"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

func TestSpinnerFrameAnimation(t *testing.T) {
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(5 * time.Millisecond)

	assert.Contains(t, w.Data, "◒ Loading")
	assert.Contains(t, w.Data, "◐ Loading")
	assert.Contains(t, w.Data, "◓ Loading")
	assert.Contains(t, w.Data, "◑ Loading")
}

func TestSpinnerDotsAnimation(t *testing.T) {
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output: w,
	})

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(4 * time.Second)

	assert.Contains(t, w.Data, "◑ Loading")
	assert.Contains(t, w.Data, "◒ Loading.")
	assert.Contains(t, w.Data, "◐ Loading..")
	assert.Contains(t, w.Data, "◓ Loading...")
}

func TestSpinnerTimerAnimation(t *testing.T) {
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output:    w,
		Indicator: prompts.SpinnerTimerIndicator,
	})

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(3 * time.Second)

	assert.Contains(t, w.Data, "◒ Loading [0s]")
	assert.Contains(t, w.Data, "◑ Loading [1s]")
	assert.Contains(t, w.Data, "◒ Loading [2s]")
}

func TestSpinnerDotsAnimationDuringCI(t *testing.T) {
	os.Setenv("CI", "true")
	defer os.Setenv("CI", "")
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(5 * time.Millisecond)

	assert.Contains(t, w.Data, "◒ Loading...")
}

func TestSpinnerRemoveDotsFromMessage(t *testing.T) {
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})

	s.Start("Loading...")
	defer s.Stop("", 0)

	time.Sleep(50 * time.Millisecond)

	assert.Contains(t, w.Data, "◒ Loading")
}

func TestSpinnerMessageMethod(t *testing.T) {
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})

	s.Start("Loading...")
	defer s.Stop("", 0)
	s.Message("Still Loading")
	time.Sleep(50 * time.Millisecond)

	assert.Contains(t, w.Data, "◒ Still Loading")
}

func TestSpinnerStopMessage(t *testing.T) {
	w := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output:        w,
		FrameInterval: time.Millisecond,
	})

	s.Start("Loading...")
	time.Sleep(2 * time.Millisecond)
	s.Stop("Loaded", 0)
	time.Sleep(1 * time.Millisecond)

	assert.Contains(t, w.Data, "◇ Loaded\n")
}
