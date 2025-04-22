package prompts

import (
	"os"
	"strings"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/third_party/picocolors"
	"golang.org/x/term"
)

type TaskLogOptions struct {
	Output    *os.File
	Title     string
	Limit     int
	Spacing   int
	RetainLog bool
}

type TaskLogMessageOptions struct {
	Raw bool
}

type TaskLogCompletionOptions struct {
	ShowLog bool
}

type TaskLogger struct {
	options           TaskLogOptions
	output            *os.File
	columns           int
	secondarySymbol   string
	spacing           int
	barSize           int
	retainLog         bool
	buffer            string
	fullBuffer        string
	lastMessageWasRaw bool
	isCI              bool
}

func TaskLog(options TaskLogOptions) *TaskLogger {
	output := options.Output
	if output == nil {
		options.Output = os.Stdout
	}

	columns, _, err := term.GetSize(int(output.Fd()))
	if err != nil {
		columns = 80 // fallback
	}

	spacing := options.Spacing
	if spacing == 0 {
		spacing = 1
	}

	logger := &TaskLogger{
		options:         options,
		output:          output,
		columns:         columns,
		secondarySymbol: picocolors.Dim("|"),
		spacing:         spacing,
		barSize:         3,
		retainLog:       options.RetainLog,
		isCI:            os.Getenv("CI") == "true",
	}

	logger.write(logger.secondarySymbol + "\n")
	logger.write(picocolors.Green("✔") + "  " + options.Title + "\n")
	for range spacing {
		logger.write(logger.secondarySymbol + "\n")
	}

	return logger
}

func (t *TaskLogger) Message(msg string, options *TaskLogMessageOptions) {
	t.clear(false)

	if (options == nil || !options.Raw || !t.lastMessageWasRaw) && t.buffer != "" {
		t.buffer += "\n"
	}
	t.buffer += msg
	t.lastMessageWasRaw = options != nil && options.Raw

	if t.options.Limit > 0 {
		lines := strings.Split(t.buffer, "\n")
		if len(lines) > t.options.Limit {
			removed := lines[:len(lines)-t.options.Limit]
			lines = lines[len(lines)-t.options.Limit:]
			if t.retainLog {
				t.fullBuffer += (t.fullBuffer + "\n") + strings.Join(removed, "\n")
			}
		}
		t.buffer = strings.Join(lines, "\n")
	}

	if !t.isCI {
		t.printBuffer(t.buffer, 0)
	}
}

func (t *TaskLogger) Error(message string, options *TaskLogCompletionOptions) {
	t.clear(true)
	Error(message)
	Error(t.output, "%s %s\n", t.secondarySymbol, "\x1b[31m"+message+"\x1b[39m")
	if options == nil || options.ShowLog {
		t.renderBuffer()
	}
	t.buffer = ""
	t.fullBuffer = ""
}

func (t *TaskLogger) Success(message string, options *TaskLogCompletionOptions) {
	t.clear(true)
	fmt.Fprintf(t.output, "%s %s\n", t.secondarySymbol, picocolors.Green(message))
	if options != nil && options.ShowLog {
		t.renderBuffer()
	}
	t.buffer = ""
	t.fullBuffer = ""
}

func (t *TaskLogger) clear(clearTitle bool) {
	if t.buffer == "" {
		return
	}
	lines := strings.Split(t.buffer, "\n")
	height := 0
	for _, line := range lines {
		if line == "" {
			height++
		} else {
			height += (utils.StrLength(line) + t.barSize) / t.columns
		}
	}
	linesToClear := height + 1
	if clearTitle {
		linesToClear += t.spacing + 2
	}
	t.write("\x1b[0G" + strings.Repeat("\x1b[1A\x1b[2K", linesToClear))
}

func (t *TaskLogger) printBuffer(buf string, spacing int) {
	for _, line := range strings.Split(buf, "\n") {
		Message(line, MessageOptions{
			Output: t.output,
			FormatLinesOptions: core.FormatLinesOptions{
				Default: core.FormatLineOptions{
					Start: t.secondarySymbol,
				},
			},
		})
	}
	for range spacing {
		t.write(picocolors.Dim(t.secondarySymbol) + "\n")
	}
}

func (t *TaskLogger) renderBuffer() {
	if t.retainLog && t.fullBuffer != "" {
		t.printBuffer(t.fullBuffer+"\n"+t.buffer, t.spacing)
	} else {
		t.printBuffer(t.buffer, t.spacing)
	}
}

func (t *TaskLogger) write(s string) {
	t.output.WriteString(s)
}
