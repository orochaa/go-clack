package prompts

// Task represents a task to be executed, with a title, a function to run, and a disabled flag.
type Task struct {
	// Title of the task to display.
	Title string
	// Function to execute the task.
	Task func(message func(msg string)) (string, error)
	// Whether the task is disabled and should be skipped.
	Disabled bool
}

// Tasks executes a list of tasks, displaying a spinner for each task and handling errors.
func Tasks(tasks []Task, options SpinnerOptions) {
	for _, task := range tasks {
		if task.Disabled {
			continue
		}
		s := Spinner(options)
		s.Start(task.Title)
		result, err := task.Task(s.Message)
		if err != nil {
			s.Stop(err.Error(), 1)
			continue
		}
		s.Stop(result, 0)
	}
}
