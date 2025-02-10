package prompts

import (
	"errors"
	"os"

	"github.com/orochaa/go-clack/core"
)

// IsCancel checks if the given error is a cancellation error (core.ErrCancelPrompt).
// It returns true if the error matches core.ErrCancelPrompt, indicating that the user
// has canceled the prompt; otherwise, it returns false.
func IsCancel(err error) bool {
	return errors.Is(err, core.ErrCancelPrompt)
}

// ExitOnError handles the termination of the program when a error occurs.
// If the error is nil, the function simply returns without taking any action.
// If the error is a cancellation error (as determined by the IsCancel function),
// the program exits with a status code of 0, indicating a successful exit without errors.
// If the error is not a cancellation error, the function logs the error message using the Error function
// and then exits the program with a status code of 1, indicating an error exit.
func ExitOnError(err error) {
	if err == nil {
		return
	}
	if IsCancel(err) {
		os.Exit(0)
		return
	}
	Error(err.Error())
	os.Exit(1)
}
