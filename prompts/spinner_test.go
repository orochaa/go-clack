package prompts_test

import (
	"os"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

const frameInterval = 80 * time.Millisecond
const dotsInterval = 8 * frameInterval

func runSpinner() (*prompts.SpinnerController, *MockWriter) {
	writer := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Output: writer,
	})
	return s, writer
}

func TestSpinnerFrameAnimation(t *testing.T) {
	s, w := runSpinner()

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(4 * frameInterval)

	assert.Contains(t, w.Data, "◒ Loading")
	assert.Contains(t, w.Data, "◐ Loading")
	assert.Contains(t, w.Data, "◓ Loading")
	assert.Contains(t, w.Data, "◑ Loading")
}

func TestSpinnerDotsAnimation(t *testing.T) {
	s, w := runSpinner()

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(4 * dotsInterval)

	assert.Contains(t, w.Data, "◒ Loading.")
	assert.Contains(t, w.Data, "◐ Loading..")
	assert.Contains(t, w.Data, "◓ Loading...")
	assert.Contains(t, w.Data, "◑ Loading")
}

func TestSpinnerDotsAnimationDuringCI(t *testing.T) {
	os.Setenv("CI", "true")
	defer os.Setenv("CI", "")
	s, w := runSpinner()

	s.Start("Loading")
	defer s.Stop("", 0)
	time.Sleep(2 * frameInterval)

	assert.Contains(t, w.Data, "◒ Loading...")
}

func TestSpinnerRemoveDotsFromMessage(t *testing.T) {
	s, w := runSpinner()

	s.Start("Loading...")
	defer s.Stop("", 0)

	time.Sleep(2 * frameInterval)

	assert.Contains(t, w.Data, "◒ Loading")
}

func TestSpinnerMessageMethod(t *testing.T) {
	s, w := runSpinner()

	s.Start("Loading...")
	defer s.Stop("", 0)
	s.Message("Still Loading")
	time.Sleep(2 * frameInterval)

	assert.Contains(t, w.Data, "◐ Still Loading")
}

func TestSpinnerStopMessage(t *testing.T) {
	s, w := runSpinner()

	s.Start("Loading...")
	time.Sleep(2 * frameInterval)
	s.Stop("Loaded", 0)
	time.Sleep(1 * frameInterval)

	assert.Contains(t, w.Data, "◇ Loaded\n")
}
