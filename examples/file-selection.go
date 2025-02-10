package main

import "github.com/orochaa/go-clack/prompts"

func FileSelection() {
	prompts.Path(prompts.PathParams{
		Message: "Input file:",
	})
	prompts.SelectPath(prompts.SelectPathParams{
		Message: "Select file:",
	})
	prompts.MultiSelectPath(prompts.MultiSelectPathParams{
		Message: "Select multiple files:",
	})
}
