package core_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/third_party/picocolors"
	"github.com/stretchr/testify/assert"
)

func newPrompt() *core.Prompt[string] {
	return core.NewPrompt(core.PromptParams[string]{
		Render: func(p *core.Prompt[string]) string { return "" },
	})
}

func TestEventEmitter(t *testing.T) {
	testCases := []struct {
		description string
		emit        core.Event
		listen      core.Event
		args        []any
		calledTimes int
	}{
		{
			description: "emit event",
			emit:        core.KeyEvent,
			listen:      core.KeyEvent,
			calledTimes: 1,
		},
		{
			description: "emit event with a single arg",
			emit:        core.KeyEvent,
			listen:      core.KeyEvent,
			args:        []any{core.Key{Name: core.RightKey}},
			calledTimes: 1,
		},
		{
			description: "emit event with a multi args",
			emit:        core.KeyEvent,
			listen:      core.KeyEvent,
			args:        []any{core.Key{Name: core.RightKey}, "value"},
			calledTimes: 1,
		},
		{
			description: "emit event multiple times",
			emit:        core.KeyEvent,
			listen:      core.KeyEvent,
			calledTimes: 3,
		},
		{
			description: "emit other event",
			emit:        core.ErrorEvent,
			listen:      core.KeyEvent,
			calledTimes: 0,
		},
	}

	for _, tC := range testCases {
		p := newPrompt()
		var calledTimes int
		p.On(tC.listen, func(args ...any) {
			calledTimes++
			assert.Equal(t, tC.args, args)
		})
		for range max(tC.calledTimes, 1) {
			p.Emit(tC.emit, tC.args...)
		}
		assert.Equal(t, tC.calledTimes, calledTimes)
	}
}

func TestEmitEventOnce(t *testing.T) {
	p := newPrompt()
	calledTimes := 0

	p.Once(core.KeyEvent, func(args ...any) {
		calledTimes++
	})
	p.Emit(core.KeyEvent)
	p.Emit(core.KeyEvent)
	assert.Equal(t, 1, calledTimes)
}

func TestEmitUnsubscribedEvent(t *testing.T) {
	p := newPrompt()
	calledTimes := 0
	listener := func(args ...any) {
		calledTimes++
	}

	p.On(core.KeyEvent, listener)
	p.Off(core.KeyEvent, listener)
	p.Emit(core.KeyEvent)
	assert.Equal(t, 0, calledTimes)
}

func TestParseKey(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, core.Key{Name: core.EnterKey}, *p.ParseKey('\n'))
	assert.Equal(t, core.Key{Name: core.EnterKey}, *p.ParseKey('\r'))
	assert.Equal(t, core.Key{Name: core.SpaceKey}, *p.ParseKey(' '))
	assert.Equal(t, core.Key{Name: core.BackspaceKey}, *p.ParseKey('\b'))
	assert.Equal(t, core.Key{Name: core.BackspaceKey}, *p.ParseKey(127))
	assert.Equal(t, core.Key{Name: core.TabKey}, *p.ParseKey('\t'))
	assert.Equal(t, core.Key{Name: core.CancelKey}, *p.ParseKey(3))
	assert.Equal(t, core.Key{Name: "a", Char: "a"}, *p.ParseKey('a'))
}

func TestTrackValue(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, "", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Char: "a"}, "", 0)
	assert.Equal(t, "a", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Char: "b"}, "a", 1)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Char: "c"}, "ab", 1)
	assert.Equal(t, "acb", p.Value)
	assert.Equal(t, 2, p.CursorIndex)
}

func TestTrackCursor(t *testing.T) {
	p := newPrompt()

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.HomeKey}, "abc", 3)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.EndKey}, "abc", 0)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.LeftKey}, "abc", 3)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.LeftKey}, "abc", 0)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.RightKey}, "abc", 2)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.RightKey}, "abc", 3)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTrackBackspace(t *testing.T) {
	p := newPrompt()

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, "abc", 3)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, "abc", 2)
	assert.Equal(t, "ac", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, "abc", 1)
	assert.Equal(t, "bc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value, p.CursorIndex = p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, "abc", 0)
	assert.Equal(t, "abc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)
}

func TestTrackState(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, core.InitialState, p.State)

	p.PressKey(&core.Key{})
	assert.Equal(t, core.ActiveState, p.State)

	p.State = core.ErrorState
	p.PressKey(&core.Key{})
	assert.Equal(t, core.ActiveState, p.State)

	p.Validate = func(value string) error { return errors.New("") }
	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)

	p.Validate = nil
	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.SubmitState, p.State)

	p.PressKey(&core.Key{Name: core.CancelKey})
	assert.Equal(t, core.CancelState, p.State)
}

func TestValidateValue(t *testing.T) {
	p := newPrompt()
	p.Validate = func(value string) error {
		return fmt.Errorf("invalid value: %v", value)
	}

	p.Value = "foo"
	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, "invalid value: foo", p.Error)
}

func TestEmitFinalizeBeforeSubmit(t *testing.T) {
	p := newPrompt()
	calledTimes := 0
	p.On(core.FinalizeEvent, func(args ...any) {
		calledTimes++
	})
	p.On(core.SubmitEvent, func(args ...any) {
		assert.Equal(t, 1, calledTimes)
		calledTimes++
	})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, 2, calledTimes)
}

func TestAsyncValidation(t *testing.T) {
	p := newPrompt()
	p.Validate = func(value string) error {
		time.Sleep(530 * time.Millisecond)
		return nil
	}

	go p.PressKey(&core.Key{Name: core.EnterKey})
	time.Sleep(400 * time.Millisecond)
	assert.Equal(t, true, p.IsValidating)
	assert.Equal(t, core.ValidateState, p.State)
	assert.GreaterOrEqual(t, p.ValidationDuration, 400*time.Millisecond)

	time.Sleep(126 * time.Millisecond)
	assert.Equal(t, true, p.IsValidating)
	assert.Equal(t, core.ValidateState, p.State)
	assert.GreaterOrEqual(t, p.ValidationDuration, 525*time.Millisecond)

	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, false, p.IsValidating)
	assert.Equal(t, core.SubmitState, p.State)
}

func TestDiffLines(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, []int{1, 2}, p.DiffLines("a", "a\nb\nc"))
	assert.Equal(t, []int{1, 2}, p.DiffLines("a\nb\nc", "a"))
	assert.Equal(t, []int{1, 2}, p.DiffLines("a\nb\nc", "a\nc\nb"))
	assert.Equal(t, []int(nil), p.DiffLines("a\nb\nc", "a\nb\nc"))
}

func TestLimitLines(t *testing.T) {
	testCases := []struct {
		description string
		usedLines   int
		frameHeight int
		cursorIndex int
		expected    string
	}{
		{
			description: "do not limit lines with frame's height <= terminal's height",
			frameHeight: 10,
			expected:    strings.Join([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, "\r\n"),
		},
		{
			description: "limit lines with frame's height > terminal's height",
			frameHeight: 11,
			expected:    strings.Join([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "..."}, "\r\n"),
		},
		{
			description: "limit lines with usedLine + frame's height > terminal's height",
			frameHeight: 10,
			usedLines:   3,
			expected:    strings.Join([]string{"1", "2", "3", "4", "5", "6", "..."}, "\r\n"),
		},
		{
			description: "limit lines with cursor on the middle start of the list",
			frameHeight: 20,
			cursorIndex: 8,
			expected:    strings.Join([]string{"...", "3", "4", "5", "6", "7", "8", "9", "10", "..."}, "\r\n"),
		},
		{
			description: "limit lines with cursor on the middle of the list",
			frameHeight: 20,
			cursorIndex: 10,
			expected:    strings.Join([]string{"...", "5", "6", "7", "8", "9", "10", "11", "12", "..."}, "\r\n"),
		},
		{
			description: "limit lines with cursor on the middle end of the list",
			frameHeight: 20,
			cursorIndex: 16,
			expected:    strings.Join([]string{"...", "11", "12", "13", "14", "15", "16", "17", "18", "..."}, "\r\n"),
		},
		{
			description: "limit lines with cursor on the end of the list",
			frameHeight: 20,
			cursorIndex: 17,
			expected:    strings.Join([]string{"...", "12", "13", "14", "15", "16", "17", "18", "19", "20"}, "\r\n"),
		},
		{
			description: "limit lines with cursor at the end of the list",
			frameHeight: 20,
			cursorIndex: 20,
			expected:    strings.Join([]string{"...", "12", "13", "14", "15", "16", "17", "18", "19", "20"}, "\r\n"),
		},
	}

	p := newPrompt()

	for _, tC := range testCases {
		lines := make([]string, tC.frameHeight)
		for i := range tC.frameHeight {
			lines[i] = fmt.Sprint(i + 1)
		}

		p.CursorIndex = tC.cursorIndex
		frame := p.LimitLines(lines, tC.usedLines)
		assert.Equal(t, tC.expected, frame)
	}
}

func TestFormatLines(t *testing.T) {
	testCases := []struct {
		description string
		lines       []string
		options     core.FormatLinesOptions
		expected    string
	}{
		{
			description: "format first line start",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Start: "-"},
			},
			expected: strings.Join([]string{"- a", "b", "c"}, "\r\n"),
		},
		{
			description: "format new line start",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				NewLine: core.FormatLineOptions{Start: "-"},
			},
			expected: strings.Join([]string{"a", "- b", "c"}, "\r\n"),
		},
		{
			description: "format last line start",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				LastLine: core.FormatLineOptions{Start: "-"},
			},
			expected: strings.Join([]string{"a", "b", "- c"}, "\r\n"),
		},
		{
			description: "format first line end",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{End: "-"},
			},
			expected: strings.Join([]string{"a -", "b", "c"}, "\r\n"),
		},
		{
			description: "format new line end",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				NewLine: core.FormatLineOptions{End: "-"},
			},
			expected: strings.Join([]string{"a", "b -", "c"}, "\r\n"),
		},
		{
			description: "format last line end",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				LastLine: core.FormatLineOptions{End: "-"},
			},
			expected: strings.Join([]string{"a", "b", "c -"}, "\r\n"),
		},
		{
			description: "format first line sides",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Sides: "-"},
			},
			expected: strings.Join([]string{"- a -", "b", "c"}, "\r\n"),
		},
		{
			description: "format new line sides",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				NewLine: core.FormatLineOptions{Sides: "-"},
			},
			expected: strings.Join([]string{"a", "- b -", "c"}, "\r\n"),
		},
		{
			description: "format last line sides",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				LastLine: core.FormatLineOptions{Sides: "-"},
			},
			expected: strings.Join([]string{"a", "b", "- c -"}, "\r\n"),
		},
		{
			description: "format first line style",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{
					Style: func(line string) string {
						return "-" + line + "-"
					},
				},
			},
			expected: strings.Join([]string{"-a-", "b", "c"}, "\r\n"),
		},
		{
			description: "format new line style",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				NewLine: core.FormatLineOptions{
					Style: func(line string) string {
						return "-" + line + "-"
					},
				},
			},
			expected: strings.Join([]string{"a", "-b-", "c"}, "\r\n"),
		},
		{
			description: "format last line style",
			lines:       []string{"a", "b", "c"},
			options: core.FormatLinesOptions{
				LastLine: core.FormatLineOptions{
					Style: func(line string) string {
						return "-" + line + "-"
					},
				},
			},
			expected: strings.Join([]string{"a", "b", "-c-"}, "\r\n"),
		},
		{
			description: "format unique line with fist and last line start options",
			lines:       []string{"a"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Sides: "-"},
				LastLine:  core.FormatLineOptions{Sides: "*"},
			},
			expected: strings.Join([]string{"* a *"}, "\r\n"),
		},
		{
			description: "format unique line with fist and last line merged options",
			lines:       []string{"a"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Start: "-"},
				LastLine:  core.FormatLineOptions{End: "*"},
			},
			expected: strings.Join([]string{"- a *"}, "\r\n"),
		},
		{
			description: "format new line with last start options",
			lines:       []string{"a", "b"},
			options: core.FormatLinesOptions{
				NewLine:  core.FormatLineOptions{Sides: "-"},
				LastLine: core.FormatLineOptions{Sides: "*"},
			},
			expected: strings.Join([]string{"a", "* b *"}, "\r\n"),
		},
		{
			description: "format new line with last merged options",
			lines:       []string{"a", "b"},
			options: core.FormatLinesOptions{
				NewLine:  core.FormatLineOptions{Start: "-"},
				LastLine: core.FormatLineOptions{End: "*"},
			},
			expected: strings.Join([]string{"a", "- b *"}, "\r\n"),
		},
		{
			description: "format with defaults",
			lines:       []string{strings.Repeat("a", 20) + " " + strings.Repeat("b", 80), "c"},
			options: core.FormatLinesOptions{
				Default: core.FormatLineOptions{
					Sides: "|",
					Style: func(line string) string {
						return fmt.Sprintf("(%s)", line)
					},
				},
			},
			expected: strings.Join([]string{
				fmt.Sprintf("| (%s %s) |", strings.Repeat("a", 20), strings.Repeat("b", 53)),
				fmt.Sprintf("| (%s) |", strings.Repeat("b", 27)),
				fmt.Sprintf("| (%s) |", strings.Repeat("c", 1)),
			}, "\r\n"),
		},
		{
			description: "format with white spaces",
			lines:       []string{"a", " b", "  c", "x", " y", "  z"},
			options: core.FormatLinesOptions{
				Default: core.FormatLineOptions{Start: "|"},
			},
			expected: strings.Join([]string{"| a", "|  b", "|   c", "| x", "|  y", "|   z"}, "\r\n"),
		},
		{
			description: "format overflowed lines",
			lines:       []string{strings.Repeat("a", 81), strings.Repeat("b", 76), strings.Repeat("c", 79)},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Sides: "-"},
				NewLine:   core.FormatLineOptions{Sides: "|"},
				LastLine:  core.FormatLineOptions{Sides: "*"},
			},
			expected: strings.Join([]string{
				fmt.Sprintf("- %s -", strings.Repeat("a", 76)),
				fmt.Sprintf("- %s -", strings.Repeat("a", 5)),
				fmt.Sprintf("| %s |", strings.Repeat("b", 76)),
				fmt.Sprintf("* %s *", strings.Repeat("c", 76)),
				fmt.Sprintf("* %s *", strings.Repeat("c", 3)),
			}, "\r\n"),
		},
		{
			description: "format double overflowed lines",
			lines:       []string{strings.Repeat("a", 180), strings.Repeat("b", 180)},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Sides: "-"},
				Default:   core.FormatLineOptions{Sides: "*"},
				LastLine:  core.FormatLineOptions{Sides: "-"},
			},
			expected: strings.Join([]string{
				fmt.Sprintf("- %s -", strings.Repeat("a", 76)),
				fmt.Sprintf("- %s -", strings.Repeat("a", 76)),
				fmt.Sprintf("- %s -", strings.Repeat("a", 28)),
				fmt.Sprintf("- %s -", strings.Repeat("b", 76)),
				fmt.Sprintf("- %s -", strings.Repeat("b", 76)),
				fmt.Sprintf("- %s -", strings.Repeat("b", 28)),
			}, "\r\n"),
		},
		{
			description: "format box",
			lines:       []string{strings.Repeat("a", 67), "b"},
			options: core.FormatLinesOptions{
				Default:  core.FormatLineOptions{Sides: "|"},
				MinWidth: 70,
				MaxWidth: 70,
			},
			expected: strings.Join([]string{
				fmt.Sprintf("| %s%s |", strings.Repeat("a", 66), strings.Repeat(" ", 0)),
				fmt.Sprintf("| %s%s |", strings.Repeat("a", 1), strings.Repeat(" ", 65)),
				fmt.Sprintf("| %s%s |", strings.Repeat("b", 1), strings.Repeat(" ", 65)),
			}, "\r\n"),
		},
		{
			description: "format with complex options",
			lines:       []string{strings.Repeat("a", 20) + " " + strings.Repeat("b", 80), "c"},
			options: core.FormatLinesOptions{
				FirstLine: core.FormatLineOptions{Start: "-"},
				NewLine:   core.FormatLineOptions{End: "-"},
				LastLine:  core.FormatLineOptions{Start: "=", End: "="},
				Default:   core.FormatLineOptions{Sides: "*"},
			},
			expected: strings.Join([]string{
				fmt.Sprintf("- %s %s *", strings.Repeat("a", 20), strings.Repeat("b", 55)),
				fmt.Sprintf("- %s *", strings.Repeat("b", 25)),
				fmt.Sprintf("= %s =", strings.Repeat("c", 1)),
			}, "\r\n"),
		},
		{
			description: "format lines with ansi colors",
			lines:       []string{picocolors.Inverse("T")},
			options: core.FormatLinesOptions{
				Default:  core.FormatLineOptions{Start: "|"},
				LastLine: core.FormatLineOptions{End: "|"},
				MinWidth: 80,
			},
			expected: fmt.Sprintf("| T%s |", strings.Repeat(" ", 75)),
		},
		{
			description: "format blank line",
			lines:       []string{"    "},
			options: core.FormatLinesOptions{
				Default:  core.FormatLineOptions{Sides: "|"},
				MinWidth: 80,
			},
			expected: fmt.Sprintf("| %s |", strings.Repeat(" ", 76)),
		},
		{
			description: "format text",
			lines:       []string{"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s"},
			options: core.FormatLinesOptions{
				Default:  core.FormatLineOptions{Start: "|"},
				MaxWidth: 50,
			},
			expected: strings.Join([]string{
				"| Lorem Ipsum is simply dummy text of the printing",
				"| and typesetting industry. Lorem Ipsum has been",
				"| the industry's standard dummy text ever since",
				"| the 1500s",
			}, "\r\n"),
		},
	}

	p := newPrompt()

	for _, tC := range testCases {
		frame := p.FormatLines(tC.lines, tC.options)
		assert.Equal(t, tC.expected, frame)
	}
}
