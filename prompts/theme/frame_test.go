package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFrame(t *testing.T) {
	frame := NewFrame()

	assert.Equal(t, "", frame.String())
}

func TestWriteLn(t *testing.T) {
	frame := NewFrame()

	frame.WriteLn("Hello", "World")
	expected := "Hello\r\nWorld\r\n"

	assert.Equal(t, expected, frame.String())
}

func TestString(t *testing.T) {
	frame := NewFrame()

	frame.WriteLn("Test")
	expected := "Test\r\n"

	assert.Equal(t, expected, frame.String())
}

func TestRemoveTrailingCRLF(t *testing.T) {
	frame := NewFrame()

	frame.WriteLn("Line 1", "Line 2")
	frame.RemoveTrailingCRLF()
	expected := "Line 1\r\nLine 2"

	assert.Equal(t, expected, frame.String())
}
