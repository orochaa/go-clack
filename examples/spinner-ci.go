/**
 * This example addresses a issue reported in GitHub Actions where `spinner` was excessively writing messages,
 * leading to confusion and cluttered output.
 * To enhance the CI workflow and provide a smoother experience,
 * the following changes have been made only for CI environment:
 * - Messages will now only be written when a `spinner` method is called and the message updated, preventing unnecessary message repetition.
 * - There will be no loading dots animation, instead it will be always `...`
 * - Instead of erase the previous message, action that is blocked during CI, it will just write a new one.
 *
 * Issue: https://github.com/natemoo-re/clack/issues/168
 */
package main

import (
	"fmt"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
)

func SpinnerCIExample() {
	prompts.Intro("Running spinner in CI environment")

	s := prompts.Spinner(prompts.SpinnerOptions{})
	total := 6000
	progress := 0
	counter := 0

	s.Start("spinner.Start")

	for progress < total {
		if progress%1000 == 0 {
			counter++
		}
		progress = min(total, progress+100)
		s.Message(fmt.Sprintf("spinner.Message [%d]", counter))

		time.Sleep(100 * time.Millisecond)
	}

	s.Stop("spinner.Stop", 0)
	prompts.Outro("Done")
}
