package main

import (
	"context"
	"time"

	"github.com/orochaa/go-clack/prompts"
)

func RaceCondition() {
	ch := make(chan string)
	defer close(ch)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
			cancel()
			ch <- "js"
		}
	}()

	go func() {
		defer cancel()
		pkg, err := prompts.Select(prompts.SelectParams[string]{
			Context: ctx,
			Message: "Pick a project type.",
			Options: []*prompts.SelectOption[string]{
				{Value: "ts", Label: "TypeScript"},
				{Value: "js", Label: "JavaScript"},
				{Value: "coffee", Label: "CoffeeScript", Hint: "oh no"},
			},
			Required: true,
		})
		if err != nil {
			ch <- err.Error()
			return
		}

		ch <- pkg
	}()

	result := <-ch
	time.Sleep(1 * time.Millisecond) // Time for the prompt to update to cancellation's frame
	println("Result:", result)
}
