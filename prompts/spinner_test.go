package prompts_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

func runSpinner() (*prompts.SpinnerController, *MockTimer, *MockWriter) {
	timer := &MockTimer{}
	writer := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
		Timer:  timer,
		Output: writer,
	})
	return s, timer, writer
}

func TestSpinnerFrameAnimation(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading")
	defer s.Stop("", 0)
	for i := 0; i < 5; i++ {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}

	assert.Equal(t, "◒ Loading", mw.Data[2])
	assert.Equal(t, "◐ Loading", mw.Data[5])
	assert.Equal(t, "◓ Loading", mw.Data[8])
	assert.Equal(t, "◑ Loading", mw.Data[11])
}

func TestSpinnerDotsAnimation(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading")
	defer s.Stop("", 0)

	for mw.Data[len(mw.Data)-1] != "◒ Loading" {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

	for mw.Data[len(mw.Data)-1] != "◒ Loading." {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading.", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

	for mw.Data[len(mw.Data)-1] != "◒ Loading.." {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading..", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

	for mw.Data[len(mw.Data)-1] != "◒ Loading..." {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading...", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

}

func TestSpinnerDotsAnimationDuringCI(t *testing.T) {
	os.Setenv("CI", "true")
	defer os.Setenv("CI", "")
	s, mt, mw := runSpinner()

	s.Start("Loading")
	defer s.Stop("", 0)

	time.Sleep(time.Microsecond)
	mt.ResolveAll()

	assert.Equal(t, "◒ Loading...", mw.Data[2], fmt.Sprint(len(mw.Data)))
}

func TestSpinnerRemoveDotsFromMessage(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading...")
	defer s.Stop("", 0)

	time.Sleep(time.Microsecond)
	mt.ResolveAll()
	time.Sleep(time.Microsecond)

	assert.Equal(t, "◒ Loading", mw.Data[2])
}

func TestSpinnerMessageMethod(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading...")
	defer s.Stop("", 0)

	time.Sleep(time.Millisecond)
	s.Message("Still Loading")
	mt.ResolveAll()
	time.Sleep(time.Millisecond)

	assert.Equal(t, "◐ Still Loading", mw.Data[5])
}

func TestSpinnerStopMessage(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading...")
	time.Sleep(time.Millisecond)
	s.Stop("Loaded", 0)
	mt.ResolveAll()
	time.Sleep(time.Millisecond)

	assert.Equal(t, "◇ Loaded\n", mw.Data[6])
}
