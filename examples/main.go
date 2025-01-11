package main

import (
	"os"

	"github.com/Mist3rBru/go-clack/prompts"
)

func HandleCancel(err error) {
	if err != nil {
		prompts.Cancel("Operation cancelled.")
		os.Exit(0)
	}
}

func main() {
	prompt, err := prompts.Select(prompts.SelectParams[string]{
		Message: "Select a example:",
		Options: []*prompts.SelectOption[string]{
			{Label: "basic"},
			{Label: "changeset"},
			{Label: "spinner"},
			{Label: "spinner-ci"},
			{Label: "async validation"},
		},
	})
	if err != nil {
		return
	}

	print("\n")
	switch prompt {
	case "basic":
		BasicExample()
	case "changeset":
		ChangesetExample()
	case "spinner":
		SpinnerExample()
	case "spinner-ci":
		os.Setenv("CI", "true")
		SpinnerCIExample()
	case "async validation":
		AsyncValidation()
	default:
		prompts.Error("example not found")
	}
}
