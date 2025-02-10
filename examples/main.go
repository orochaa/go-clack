package main

import (
	"os"

	"github.com/orochaa/go-clack/prompts"
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
			{Label: "spinner-timer"},
			{Label: "spinner-ci"},
			{Label: "async-validation"},
			{Label: "file-selection"},
			{Label: "race-condition"},
			{Label: "custom-keys"},
		},
	})
	if err != nil {
		return
	}

	switch prompt {
	case "basic":
		BasicExample()
	case "changeset":
		ChangesetExample()
	case "spinner":
		SpinnerExample()
	case "spinner-timer":
		SpinnerTimerExample()
	case "spinner-ci":
		os.Setenv("CI", "true")
		SpinnerCIExample()
	case "async-validation":
		AsyncValidation()
	case "file-selection":
		FileSelection()
	case "race-condition":
		RaceCondition()
	case "custom-keys":
		CustomKeys()
	default:
		prompts.Error("example not found")
	}
}
