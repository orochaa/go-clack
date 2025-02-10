package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/orochaa/go-clack/prompts"
	"github.com/orochaa/go-clack/third_party/picocolors"
	"github.com/orochaa/go-clack/third_party/sisteransi"
)

func BasicExample() {
	os.Stdout.Write([]byte(sisteransi.EraseDown()))
	prompts.Intro(picocolors.BgCyan(picocolors.Black(" create-app ")))

	path, err := prompts.Path(prompts.PathParams{
		Message:      "Where should we create your project?",
		InitialValue: ".",
		Validate: func(value string) error {
			if value == "" {
				return errors.New("please enter a path")
			}
			if !strings.HasPrefix(value, ".") {
				return errors.New("please enter a relative path")
			}
			return nil
		},
	})
	HandleCancel(err)

	_, err = prompts.Password(prompts.PasswordParams{
		Message: "Provide a password",
		Validate: func(value string) error {
			if value == "" {
				return errors.New("please enter a password")
			}
			if len(value) < 5 {
				return errors.New("password should have at least 5 characters")
			}
			return nil
		},
	})
	HandleCancel(err)

	_, err = prompts.Select(prompts.SelectParams[string]{
		Message:      fmt.Sprint("Pick a project type within", path, "."),
		InitialValue: "ts",
		Options: []*prompts.SelectOption[string]{
			{Label: "TypeScript", Value: "ts"},
			{Label: "JavaScript", Value: "js"},
			{Label: "Go", Value: "go"},
			{Label: "Python", Value: "python"},
			{Label: "CoffeeScript", Value: "coffee", Hint: "oh no"},
			{Label: "Rust", Value: "rust"},
		},
	})
	HandleCancel(err)

	_, err = prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message:      "Select additional tools.",
		InitialValue: []string{"prettier", "eslint"},
		Options: []*prompts.MultiSelectOption[string]{
			{Label: "Prettier", Value: "prettier", Hint: "recommended"},
			{Label: "ESLint", Value: "eslint", Hint: "recommended"},
			{Label: "Stylelint", Value: "stylelint"},
			{Label: "GitHub Action", Value: "gh-action"},
		},
	})
	HandleCancel(err)

	install, err := prompts.Confirm(prompts.ConfirmParams{
		Message:      "Install dependencies?",
		InitialValue: false,
	})
	HandleCancel(err)

	if install {
		s := prompts.Spinner(prompts.SpinnerOptions{})
		s.Start("Installing via pnpm")
		time.Sleep(3 * time.Second)
		s.Stop("Installed via pnpm", 0)
	}

	var installMsg string
	if install {
		installMsg = "pnpm install\npnpm dev"
	}

	nextSteps := fmt.Sprintf("cd %s\n%s", path, installMsg)

	prompts.Note(nextSteps, prompts.NoteOptions{Title: "Next steps."})

	prompts.Outro(fmt.Sprintf("Problems? %s", picocolors.Underline(picocolors.Cyan("https://example.com/issues"))))
}
