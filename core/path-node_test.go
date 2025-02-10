package core_test

import (
	"os"
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPathNode(rootPath string) *core.PathNode {
	return core.NewPathNode(rootPath, core.PathNodeOptions{FileSystem: MockFileSystem{}})
}

func TestPathNodeRoot(t *testing.T) {
	cwd, _ := os.Getwd()
	node := core.NewPathNode(cwd, core.PathNodeOptions{
		OnlyShowDir: true,
		FileSystem:  MockFileSystem{},
	})

	assert.Equal(t, cwd, node.Path)
	assert.Equal(t, cwd, node.Name)
	assert.Equal(t, true, node.IsDir)
	assert.Equal(t, 0, node.Index)
	assert.Equal(t, 0, node.Depth)
	assert.Equal(t, (*core.PathNode)(nil), node.Parent)
	assert.Equal(t, true, node.OnlyShowDir)
	assert.Greater(t, len(node.Children), 0)
}

func TestPathNodeIsRoot(t *testing.T) {
	node := newPathNode("/root/go-clack/core")
	other := newPathNode("/root/go-clack/prompts")

	assert.True(t, node.IsRoot())

	node.Parent = other
	assert.False(t, node.IsRoot())
}

func TestPathNodeIsEqual(t *testing.T) {
	node := newPathNode("/root/go-clack/core")
	other := newPathNode("/root/go-clack/prompts")

	assert.True(t, node.IsEqual(node))
	assert.False(t, node.IsEqual(other))
}

func TestPathNodeTraverseNodes(t *testing.T) {
	node := newPathNode("/root/go-clack/core")
	node.Open()

	counter := 0
	node.TraverseNodes(func(node *core.PathNode) {
		counter++
	})

	assert.Equal(t, len(node.Children)+1, counter)
}

func TestPathNodeFlat(t *testing.T) {
	node := newPathNode("/root/go-clack/core")

	assert.Equal(t, 3, len(node.Flat()))
}

func TestPathNodeFilteredFlat(t *testing.T) {
	root := newPathNode("/root/go-clack/core")
	var options []*core.PathNode

	options = root.FilteredFlat("", root)
	assert.Equal(t, 3, len(options))

	options = root.FilteredFlat("d", root)
	assert.Equal(t, 3, len(options))

	options = root.FilteredFlat("", root.Children[0])
	assert.Equal(t, 3, len(options))

	options = root.FilteredFlat("d", root.Children[0])
	assert.Equal(t, 2, len(options))

	options = root.FilteredFlat("f", root.Children[0])
	assert.Equal(t, 2, len(options))
}

func TestPathNodeIndexOf(t *testing.T) {
	node := newPathNode("/root/go-clack/core")
	node.Open()

	assert.Equal(t, 1, node.IndexOf(node.Children[1], node.Children))
}
