package core

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/orochaa/go-clack/core/internals"
)

type PathNode struct {
	Index  int
	Depth  int
	Path   string
	Name   string
	Parent *PathNode

	IsDir    bool
	IsOpen   bool
	Children []*PathNode

	IsSelected bool

	FileSystem  FileSystem
	OnlyShowDir bool
}

func (n *PathNode) String() string {
	if n == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"{Index:%d, Depth:%d, Name:%q, IsDir:%t, IsOpen:%t, Children:%d, IsSelected:%t}",
		n.Index,
		n.Depth,
		n.Name,
		n.IsDir,
		n.IsOpen,
		len(n.Children),
		n.IsSelected,
	)
}

type PathNodeOptions struct {
	OnlyShowDir bool
	FileSystem  FileSystem
}

// NewPathNode initializes a new PathNode with the provided root path and options.
// It sets up the root node, configures the file system, and opens the node to populate its children.
//
// Parameters:
//   - rootPath (string): The root path for the node.
//   - options (PathNodeOptions): Configuration options for the node.
//   - OnlyShowDir (bool): Whether to only show directories (default: false).
//   - FileSystem (FileSystem): The file system implementation to use (default: OSFileSystem).
//
// Returns:
//   - *PathNode: A new instance of PathNode.
func NewPathNode(rootPath string, options PathNodeOptions) *PathNode {
	if options.FileSystem == nil {
		options.FileSystem = internals.OSFileSystem{}
	}

	root := &PathNode{
		Path:  rootPath,
		Name:  rootPath,
		IsDir: true,

		OnlyShowDir: options.OnlyShowDir,
		FileSystem:  options.FileSystem,
	}
	root.Open()

	return root
}

// Open opens the current PathNode, reads its directory entries, and populates its children.
// If the node is not a directory or is already open, this function does nothing.
func (p *PathNode) Open() {
	if !p.IsDir || p.IsOpen {
		return
	}

	entries, err := p.FileSystem.ReadDir(p.Path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if p.OnlyShowDir && !entry.IsDir() {
			continue
		}
		p.Children = append(p.Children, &PathNode{
			Depth:  p.Depth + 1,
			Path:   path.Join(p.Path, entry.Name()),
			Name:   entry.Name(),
			Parent: p,
			IsDir:  entry.IsDir(),

			FileSystem:  p.FileSystem,
			OnlyShowDir: p.OnlyShowDir,
		})
	}

	sort.SliceStable(p.Children, func(i, j int) bool {
		if p.Children[i].IsDir != p.Children[j].IsDir {
			return p.Children[i].IsDir
		}
		return strings.ToLower(p.Children[i].Name) < strings.ToLower(p.Children[j].Name)
	})

	for i, child := range p.Children {
		child.Index = i
	}

	p.IsOpen = true
}

// Close closes the current PathNode by clearing its children and marking it as closed.
func (p *PathNode) Close() {
	p.Children = []*PathNode(nil)
	p.IsOpen = false
}

// TraverseNodes traverses the node and its children, applying the provided visit function to each node.
//
// Parameters:
//   - visit (func(node *PathNode)): A function to apply to each node during traversal.
func (p *PathNode) TraverseNodes(visit func(node *PathNode)) {
	var traverse func(node *PathNode)
	traverse = func(node *PathNode) {
		visit(node)
		if !node.IsDir {
			return
		}
		for _, child := range node.Children {
			traverse(child)
		}
	}

	traverse(p)
}

// Flat returns a flattened list of all nodes in the tree, starting from the current node.
//
// Returns:
//   - []*PathNode: A slice of all nodes in the tree.
func (p *PathNode) Flat() []*PathNode {
	var options []*PathNode
	p.TraverseNodes(func(node *PathNode) {
		options = append(options, node)
	})
	return options
}

// FilteredFlat returns a filtered and flattened list of nodes based on the provided search term.
// If the search term is empty or invalid, it returns the full flattened list.
//
// Parameters:
//   - search (string): The search term to filter nodes by.
//   - currentNode (*PathNode): The current node to filter relative to.
//
// Returns:
//   - []*PathNode: A slice of filtered nodes.
func (p *PathNode) FilteredFlat(search string, currentNode *PathNode) []*PathNode {
	searchRegex, err := regexp.Compile("(?i)" + search)
	if err != nil || search == "" {
		return p.Flat()
	}

	var options []*PathNode
	p.TraverseNodes(func(node *PathNode) {
		if node.Depth == currentNode.Depth && node.Depth > 0 {
			if matched := searchRegex.MatchString(node.Name); matched {
				options = append(options, node)
			}
		} else {
			options = append(options, node)
		}
	})

	return options
}

// Layer returns the children of the current node's parent, representing the current layer in the tree.
//
// Returns:
//   - []*PathNode: A slice of nodes in the current layer.
func (p *PathNode) Layer() []*PathNode {
	if p.IsRoot() {
		return []*PathNode{p}
	}

	return p.Parent.Children
}

// FilteredLayer returns a filtered list of nodes in the current layer based on the provided search term.
// If the search term is empty or invalid, it returns the full layer.
//
// Parameters:
//   - search (string): The search term to filter nodes by.
//
// Returns:
//   - []*PathNode: A slice of filtered nodes in the current layer.
func (p *PathNode) FilteredLayer(search string) []*PathNode {
	searchRegex, err := regexp.Compile("(?i)" + search)
	if err != nil || search == "" {
		return p.Layer()
	}

	var layer []*PathNode
	for _, node := range p.Layer() {
		if matched := searchRegex.MatchString(node.Name); matched {
			layer = append(layer, node)
		}
	}

	return layer
}

// FirstChild returns the first child of the current node.
// If the node has no children, it returns nil.
//
// Returns:
//   - *PathNode: The first child node.
func (p *PathNode) FirstChild() *PathNode {
	if len(p.Children) == 0 {
		return nil
	}
	return p.Children[0]
}

// LastChild returns the last child of the current node.
// If the node has no children, it returns nil.
//
// Returns:
//   - *PathNode: The last child node.
func (p *PathNode) LastChild() *PathNode {
	if len(p.Children) == 0 {
		return nil
	}
	return p.Children[len(p.Children)-1]
}

// PrevChild returns the previous child relative to the provided index.
// If the index is out of bounds, it wraps around to the last child.
//
// Parameters:
//   - index (int): The index of the current child.
//
// Returns:
//   - *PathNode: The previous child node.
func (p *PathNode) PrevChild(index int) *PathNode {
	if index <= 0 {
		return p.LastChild()
	}
	return p.Children[index-1]
}

// NextChild returns the next child relative to the provided index.
// If the index is out of bounds, it wraps around to the first child.
//
// Parameters:
//   - index (int): The index of the current child.
//
// Returns:
//   - *PathNode: The next child node
func (p *PathNode) NextChild(index int) *PathNode {
	if index+1 >= len(p.Children) {
		return p.FirstChild()
	}
	return p.Children[index+1]
}

// IsRoot checks if the current node is the root node.
// It returns true if the node has no parent, false otherwise.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (p *PathNode) IsRoot() bool {
	return p.Parent == nil
}

// IsEqual checks if the current node is equal to another node based on their paths.
//
// Parameters:
//   - node (*PathNode): The node to compare with.
//
// Returns:
//   - bool: True if the nodes are equal, false otherwise.
func (p *PathNode) IsEqual(node *PathNode) bool {
	return node.Path == p.Path
}

// IndexOf returns the index of a given node in the provided options slice.
// If the node is not found, it returns -1.
//
// Parameters:
//   - node (*PathNode): The node to find the index of.
//   - options ([]*PathNode): The slice of nodes to search in.
//
// Returns:
//   - int: The index of the node, or -1 if not found.
func (p *PathNode) IndexOf(node *PathNode, options []*PathNode) int {
	if node != nil {
		for i, option := range options {
			if option.IsEqual(node) {
				return i
			}
		}
	}
	return -1
}
