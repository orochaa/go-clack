package core

import "bytes"

// Frame represents a container for building and managing text with line-by-line operations.
type Frame struct {
	bytes.Buffer
}

// NewFrame creates and returns a new Frame instance.
func NewFrame() Frame {
	return Frame{*new(bytes.Buffer)}
}

// WriteLn writes one or more lines to the buffer, appending "\r\n" to each line.
func (f *Frame) WriteLn(lines ...string) {
	for _, line := range lines {
		f.WriteString(line + "\r\n")
	}
}

// RemoveTrailingCRLF removes the trailing "\r\n" from the buffer if it exists.
func (f *Frame) RemoveTrailingCRLF() {
	data := f.Bytes()
	if len(data) >= 2 && string(data[len(data)-2:]) == "\r\n" {
		f.Truncate(len(data) - 2)
	}
}
