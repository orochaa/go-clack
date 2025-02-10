package main

import (
	"time"

	"github.com/orochaa/go-clack/prompts"
)

func SpinnerTimerExample() {
	prompts.Intro("spinner start...")

	s := prompts.Spinner(prompts.SpinnerOptions{Indicator: prompts.SpinnerTimerIndicator})

	s.Start("First spinner")
	time.Sleep(3 * time.Second)
	s.Stop("Done first spinner", 0)

	s.Start("Second spinner")
	time.Sleep(5 * time.Second)
	s.Stop("Done second spinner", 0)

	prompts.Outro("spinner stop.")
}
