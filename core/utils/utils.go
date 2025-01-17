package utils

func isControlCharacter(r rune) bool {
	return r <= 0x1f || (r >= 0x7f && r <= 0x9f)
}

func isCombiningCharacter(r rune) bool {
	return r >= 0x300 && r <= 0x36f
}

func isSurrogatePair(r rune) bool {
	return r >= 0xd800 && r <= 0xdbff
}

func StrLength(str string) int {
	if len(str) == 0 {
		return 0
	}

	length := 0
	inEscapeCode := false

	for i := 0; i < len(str); i++ {
		r := rune(str[i])

		if inEscapeCode {
			if r == 'm' {
				inEscapeCode = false
			}
			continue
		}

		if r == '\x1b' {
			inEscapeCode = true
			// length++ // count the escape code as 1 character
			continue
		}

		if isControlCharacter(r) || isCombiningCharacter(r) {
			continue
		}

		if isSurrogatePair(r) {
			i++
		}

		length++
	}

	return length
}

func MinMaxIndex(index int, max int) int {
	if index < 0 {
		return max - 1
	}
	if index >= max {
		return 0
	}
	return index
}

// SplitLines splits a string into a slice of lines.
// It handles both "\n" and "\r\n" line endings.
// If the string doesn't end with a newline, the last line is still appended to the result.
func SplitLines(str string) []string {
	if len(str) == 0 {
		return []string{""}
	}

	var lines []string
	start := 0

	for i := 0; i < len(str); i++ {
		if str[i] == '\r' {
			lines = append(lines, str[start:i])
			if i+1 < len(str) && str[i+1] == '\n' {
				// Skip next \n
				i++
			}
			start = i + 1
		} else if str[i] == '\n' {
			lines = append(lines, str[start:i])
			start = i + 1
		}
	}

	if start < len(str) {
		lines = append(lines, str[start:])
	}

	lastChar := str[len(str)-1]
	lastLine := lines[len(lines)-1]
	if lastLine != "" && (lastChar == '\r' || lastChar == '\n') {
		lines = append(lines, "")
	}

	return lines
}
