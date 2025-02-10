package utils_test

import (
	"testing"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/third_party/picocolors"

	"github.com/stretchr/testify/assert"
)

func TestStrLength(t *testing.T) {
	assert.Equal(t, 1, utils.StrLength(picocolors.Inverse(" ")))
	assert.Equal(t, 1, utils.StrLength(picocolors.Inverse("█")))
	assert.Equal(t, 5, utils.StrLength(picocolors.Cyan("| foo")))
	assert.Equal(t, 5, utils.StrLength(picocolors.Gray("|")+" "+picocolors.Dim("foo")))
	assert.Equal(t, 1, utils.StrLength(picocolors.Green("◆")))
	assert.Equal(t, 5, utils.StrLength(picocolors.Green("◇")+" "+"Foo"))
	assert.Equal(t, 5, utils.StrLength(picocolors.Green("o")+" "+"Foo"))
}

func TestSplitLines(t *testing.T) {
	assert.Equal(t, []string{""}, utils.SplitLines(""), `""`)
	assert.Equal(t, []string{"Hello, World!"}, utils.SplitLines("Hello, World!"), `"Hello, World!"`)
	assert.Equal(t, []string{"Hello, World!", ""}, utils.SplitLines("Hello, World!\n"), `"Hello, World!\n"`)
	assert.Equal(t, []string{"", "Hello, World!"}, utils.SplitLines("\nHello, World!"), `"\nHello, World!"`)
	assert.Equal(t, []string{"", "Hello, World!", ""}, utils.SplitLines("\nHello, World!\n"), `"\nHello, World!\n"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3"}, utils.SplitLines("Line 1\nLine 2\nLine 3"), `"Line 1\nLine 2\nLine 3"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3"}, utils.SplitLines("Line 1\r\nLine 2\r\nLine 3"), `"Line 1\r\nLine 2\r\nLine 3"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3"}, utils.SplitLines("Line 1\r\nLine 2\nLine 3"), `"Line 1\r\nLine 2\nLine 3"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3", ""}, utils.SplitLines("Line 1\r\nLine 2\r\nLine 3\n"), `"Line 1\r\nLine 2\r\nLine 3\n"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3", ""}, utils.SplitLines("Line 1\r\nLine 2\nLine 3\r\n"), `"Line 1\r\nLine 2\nLine 3\r\n"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3", ""}, utils.SplitLines("Line 1\nLine 2\r\nLine 3\r\n"), `"Line 1\nLine 2\r\nLine 3\r\n"`)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3", ""}, utils.SplitLines("Line 1\rLine 2\rLine 3\r"), `"Line 1\rLine 2\rLine 3\r"`)
	assert.Equal(t, []string{"", "Line 1", "Line 2", "Line 3", ""}, utils.SplitLines("\r\nLine 1\nLine 2\nLine 3\r\n"), `"\r\nLine 1\nLine 2\nLine 3\r\n"`)
	assert.Equal(t, []string{"", "Line 1", "", "Line 2", "", "Line 3", ""}, utils.SplitLines("\r\nLine 1\n\nLine 2\n\nLine 3\r\n"), `"\r\nLine 1\n\nLine 2\n\nLine 3\r\n"`)
	assert.Equal(t, []string{"", "", ""}, utils.SplitLines("\n\n\n"), `"\n\n\n"`)
	assert.Equal(t, []string{"", "", ""}, utils.SplitLines("\r\n\r\n\r\n"), `"\r\n\r\n\r\n"`)
	assert.Equal(t, []string{"", "", ""}, utils.SplitLines("\n\r\n\r\n"), `"\n\r\n\r\n"`)
	assert.Equal(t, []string{"", "", ""}, utils.SplitLines("\r\n\n\r\n"), `"\r\n\n\r\n"`)
	assert.Equal(t, []string{"", "", ""}, utils.SplitLines("\r\n\r\n\n"), `"\r\n\r\n\n"`)
}
