package utils_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"

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
