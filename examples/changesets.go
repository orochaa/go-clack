package main

import (
	"fmt"
	"strings"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

func ChangesetExample() {
	prompts.Intro(picocolors.BgCyan(picocolors.Black(" changesets ")))

	type Changeset struct {
		Packages []string
		Major    []string
		Minor    []string
		Patch    []string
		Summary  string
	}
	var changeset Changeset

	err := prompts.Workflow(&changeset).
		Step("packages", func() (any, error) {
			return prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
				Message: "Which packages would you like to include?",
				Options: map[string][]prompts.MultiSelectOption[string]{
					"changed packages": {
						{Label: "@scope/a"},
						{Label: "@scope/b"},
						{Label: "@scope/c"},
					},
					"unchanged packages": {
						{Label: "@scope/x"},
						{Label: "@scope/y"},
						{Label: "@scope/z"},
					},
				},
				Required: true,
			})
		}).
		Step("major", func() (any, error) {
			majorOptions := make([]*prompts.MultiSelectOption[string], len(changeset.Packages))
			for i, packageOption := range changeset.Packages {
				majorOptions[i] = &prompts.MultiSelectOption[string]{
					Label: packageOption,
				}
			}
			return prompts.MultiSelect(prompts.MultiSelectParams[string]{
				Message: fmt.Sprintf("Which packages should have a %s bump?", picocolors.Red("major")),
				Options: majorOptions,
			})
		}).
		Step("minor", func() (any, error) {
			minorOptions := []*prompts.MultiSelectOption[string]{}
		packagesLoop:
			for _, packageOption := range changeset.Packages {
				for _, majorPackage := range changeset.Major {
					if majorPackage == packageOption {
						continue packagesLoop
					}
				}
				minorOptions = append(minorOptions, &prompts.MultiSelectOption[string]{
					Label: packageOption,
				})
			}

			return prompts.MultiSelect(prompts.MultiSelectParams[string]{
				Message: fmt.Sprintf("Which packages should have a %s bump?", picocolors.Yellow("minor")),
				Options: minorOptions,
			})
		}).
		Step("patch", func() (any, error) {
		packagesLoop:
			for _, _package := range changeset.Packages {
				for _, majorPackage := range changeset.Major {
					if majorPackage == _package {
						continue packagesLoop
					}
				}
				for _, minorPackage := range changeset.Minor {
					if minorPackage == _package {
						continue packagesLoop
					}
				}
				changeset.Patch = append(changeset.Patch, _package)
			}

			if len(changeset.Patch) > 0 {
				note := strings.Join(changeset.Patch, picocolors.Dim(", "))
				prompts.Step(fmt.Sprintf("These packages will have a %s bump: %s", picocolors.Green("patch"), picocolors.Dim(note)))
			}

			return changeset.Patch, nil
		}).
		Step("summary", func() (any, error) {
			return prompts.Text(prompts.TextParams{
				Message:     "Please enter a summary for this change",
				Placeholder: "Summary",
			})
		}).
		Run()
	HandleCancel(err)

	prompts.Outro(fmt.Sprintf("Changeset added! %s", picocolors.Underline(picocolors.Cyan(".changeset/orange-crabs-sing.md"))))
}
