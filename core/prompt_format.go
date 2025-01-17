package core

import (
	"math"
	"strings"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

// TrackKeyValue updates the string value and cursor position based on key presses.
func (p *Prompt[TValue]) TrackKeyValue(key *Key, value string, cursorIndex int) (newValue string, newCursorIndex int) {
	switch key.Name {
	case BackspaceKey:
		if cursorIndex == 0 || len(value) == 0 {
			return value, cursorIndex
		}

		cursorIndex--
		if cursorIndex == len(value)-1 {
			return value[0:cursorIndex], cursorIndex
		}

		return value[0:cursorIndex] + value[cursorIndex+1:], cursorIndex
	case HomeKey:
		return value, 0
	case EndKey:
		return value, len(value)
	case LeftKey:
		return value, max(cursorIndex-1, 0)
	case RightKey:
		return value, min(cursorIndex+1, len(value))
	case SpaceKey:
		return value[0:cursorIndex] + " " + value[cursorIndex:], cursorIndex + 1
	}

	if len(key.Char) == 1 {
		return value[0:cursorIndex] + key.Char + value[cursorIndex:], cursorIndex + 1
	}

	return value, cursorIndex
}

// LimitLines limits the number of lines to fit within the terminal size.
func (p *Prompt[TValue]) LimitLines(lines []string, usedLines int) string {
	_, maxRows, err := p.Size()
	if err != nil {
		maxRows = 10
	}
	maxItems := max(min(maxRows-usedLines, len(lines)), 0)

	slidingWindowLocation := 0
	if p.CursorIndex >= maxItems-3 {
		slidingWindowLocation = max(min(p.CursorIndex-maxItems+3, len(lines)-maxItems), 0)
	} else if p.CursorIndex < 2 {
		slidingWindowLocation = max(p.CursorIndex-2, 0)
	}

	result := make([]string, 0, maxItems)
	shouldRenderTopEllipsis := maxItems < len(lines) && slidingWindowLocation > 0
	shouldRenderBottomEllipsis := maxItems < len(lines) && slidingWindowLocation+maxItems < len(lines)

	for i, line := range lines[slidingWindowLocation : slidingWindowLocation+maxItems] {
		isTopLimit := i == 0 && shouldRenderTopEllipsis
		isBottomLimit := i == maxItems-1 && shouldRenderBottomEllipsis
		if isTopLimit || isBottomLimit {
			result = append(result, picocolors.Dim("..."))
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\r\n")
}

type LineOption int

const (
	FirstLine LineOption = iota
	NewLine
	LastLine
)

type FormatLineOptions struct {
	Start string
	End   string
	Sides string
	Style func(line string) string
}

type FormatLinesOptions struct {
	FirstLine FormatLineOptions
	NewLine   FormatLineOptions
	LastLine  FormatLineOptions
	Default   FormatLineOptions
	MinWidth  int
	MaxWidth  int
}

// getOptionOrDefault retrieves the option for the given line and option type, or returns a default value.
func getOptionOrDefault(line LineOption, option string, options FormatLinesOptions) string {
	switch line {
	case FirstLine:
		switch option {
		case "start":
			if options.FirstLine.Start != "" {
				return options.FirstLine.Start
			} else if options.FirstLine.Sides != "" {
				return options.FirstLine.Sides
			} else if options.Default.Start != "" {
				return options.Default.Start
			}
		case "end":
			if options.FirstLine.End != "" {
				return options.FirstLine.End
			} else if options.FirstLine.Sides != "" {
				return options.FirstLine.Sides
			} else if options.Default.End != "" {
				return options.Default.End
			}
		}
	case NewLine:
		switch option {
		case "start":
			if options.NewLine.Start != "" {
				return options.NewLine.Start
			} else if options.NewLine.Sides != "" {
				return options.NewLine.Sides
			} else if options.Default.Start != "" {
				return options.Default.Start
			}
		case "end":
			if options.NewLine.End != "" {
				return options.NewLine.End
			} else if options.NewLine.Sides != "" {
				return options.NewLine.Sides
			} else if options.Default.End != "" {
				return options.Default.End
			}
		}
	case LastLine:
		switch option {
		case "start":
			if options.LastLine.Start != "" {
				return options.LastLine.Start
			} else if options.LastLine.Sides != "" {
				return options.LastLine.Sides
			} else if options.Default.Start != "" {
				return options.Default.Start
			}
		case "end":
			if options.LastLine.End != "" {
				return options.LastLine.End
			} else if options.LastLine.Sides != "" {
				return options.LastLine.Sides
			} else if options.Default.End != "" {
				return options.Default.End
			}
		}
	}
	if options.Default.Sides != "" {
		return options.Default.Sides
	}
	return ""
}

func getStyleOrDefault(line LineOption, options FormatLinesOptions) func(line string) string {
	switch line {
	case FirstLine:
		if options.FirstLine.Style != nil {
			return options.FirstLine.Style
		}
	case NewLine:
		if options.NewLine.Style != nil {
			return options.NewLine.Style
		}
	case LastLine:
		if options.LastLine.Style != nil {
			return options.LastLine.Style
		}
	}
	if options.Default.Style != nil {
		return options.Default.Style
	}
	return func(line string) string { return line }
}

func getLineOptions(options FormatLinesOptions, line LineOption) FormatLineOptions {
	return FormatLineOptions{
		Start: getOptionOrDefault(line, "start", options),
		End:   getOptionOrDefault(line, "end", options),
		Style: getStyleOrDefault(line, options),
	}
}

func mergeOptions(primary, secondary FormatLineOptions) FormatLineOptions {
	ifEmpty := func(a, b string) string {
		if a == "" {
			return b
		}
		return a
	}

	ifNil := func(a, b func(line string) string) func(line string) string {
		if a == nil {
			return b
		}
		return a
	}

	return FormatLineOptions{
		Start: ifEmpty(primary.Start, secondary.Start),
		End:   ifEmpty(primary.End, secondary.End),
		Style: ifNil(primary.Style, secondary.Style),
	}
}

// FormatLines applies styles to multiple lines based on their type and the provided options.
func (p *Prompt[TValue]) FormatLines(lines []string, options FormatLinesOptions) string {
	terminalWidth, _, err := p.Size()
	if err != nil {
		terminalWidth = 80
	}
	if options.MaxWidth == 0 {
		options.MaxWidth = math.MaxInt
	}
	minWidth := max(options.MinWidth, 0)
	maxWith := min(options.MaxWidth, terminalWidth)

	firstLine := getLineOptions(options, FirstLine)
	newLine := getLineOptions(options, NewLine)
	lastLine := getLineOptions(options, LastLine)

	formattedLines := make([]string, 0, len(lines))
	for i, line := range lines {
		var opts FormatLineOptions
		if i == 0 && len(lines) == 1 {
			opts = mergeOptions(lastLine, firstLine)
		} else if i == 0 {
			opts = firstLine
		} else if i == 1 && len(lines) == 2 {
			opts = mergeOptions(lastLine, newLine)
		} else if i+1 == len(lines) {
			opts = lastLine
		} else {
			opts = newLine
		}

		// If empty slot > 0. Sum empty space
		startEmptySlots := utils.StrLength(opts.Start)
		if startEmptySlots > 0 {
			startEmptySlots++
		}
		endEmptySlots := utils.StrLength(opts.End)
		if endEmptySlots > 0 {
			endEmptySlots++
		}
		emptySlots := startEmptySlots + endEmptySlots + utils.StrLength(opts.Style(""))

		formatAndAddLine := func(line string) {
			var startSpace, endSpace string
			if startEmptySlots > 0 {
				startSpace = " "
			}
			styledLine := opts.Style(line)
			if minWidth > 0 {
				endSpace = strings.Repeat(" ", max(minWidth+1-startEmptySlots-utils.StrLength(styledLine)-endEmptySlots, 0))
			} else if endEmptySlots > 0 {
				endSpace = " "
			}
			formattedLine := strings.Join([]string{opts.Start, startSpace, styledLine, endSpace, opts.End}, "")
			formattedLines = append(formattedLines, formattedLine)
		}

		if len(lines) == 1 && strings.Trim(lines[0], " ") == "" {
			formatAndAddLine("")
			break
		}

		var currentLine string
		for _, word := range strings.Split(line, " ") {
			if word == "" {
				currentLine += " "
			} else if strings.Trim(currentLine, " ") == "" && utils.StrLength(currentLine+word)+emptySlots <= maxWith {
				currentLine += word
			} else if utils.StrLength(currentLine+word)+emptySlots+1 <= maxWith {
				currentLine += " " + word
			} else if utils.StrLength(word)+emptySlots >= maxWith {
				var splitIndex int
				if utils.StrLength(currentLine) == 0 {
					splitIndex = maxWith - emptySlots
					formatAndAddLine(word[0:splitIndex])
				} else {
					splitIndex = maxWith - utils.StrLength(currentLine) - emptySlots - 1
					formatAndAddLine(currentLine + " " + word[0:splitIndex])
				}

				chunkLength := maxWith - emptySlots
				chunk := word[splitIndex:]
				for utils.StrLength(chunk) > chunkLength {
					formatAndAddLine(chunk[0:chunkLength])
					chunk = chunk[chunkLength:]
				}

				currentLine = chunk
			} else {
				formatAndAddLine(currentLine)
				currentLine = word
			}
		}
		formatAndAddLine(currentLine)
	}

	return strings.Join(formattedLines, "\r\n")
}
