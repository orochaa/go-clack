package main

import (
	"errors"
	"time"

	"github.com/orochaa/go-clack/prompts"
)

func AsyncValidation() {
	prompts.Text(prompts.TextParams{
		Message: "What is your GitHub's username?",
		Validate: func(value string) error {
			time.Sleep(7 * time.Second)
			return errors.New("invalid username, please try again")
		},
	})
}
