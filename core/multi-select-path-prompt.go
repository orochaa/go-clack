package core

import (
	"context"
	"os"
	"path"
	"sort"

	"github.com/Mist3rBru/go-clack/core/internals"
	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type MultiSelectPathPrompt struct {
	Prompt[[]string]
	Root          *PathNode
	CurrentOption *PathNode
	OnlyShowDir   bool
	Filter        bool
	Search        string
	Required      bool
	FileSystem    FileSystem
}

type MultiSelectPathPromptParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue []string
	InitialPath  string
	OnlyShowDir  bool
	Required     bool
	Filter       bool
	FileSystem   FileSystem
	Validate     func(value []string) error
	Render       func(p *MultiSelectPathPrompt) string
}

// NewMultiSelectPathPrompt initializes and returns a new instance of MultiSelectPathPrompt.
//
// The user can navigate through directories and files using arrow keys.
// The user can select multiple options using space key.
// The prompt returns the path of the selected options.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue ([]string): Initial selected paths (default: nil).
//   - InitialPath (string): The initial directory path to start from (default: current working directory).
//   - OnlyShowDir (bool): Whether to only show directories (default: false).
//   - Required (bool): Whether at least one option must be selected (default: false).
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - FileSystem (FileSystem): The file system implementation to use (default: OSFileSystem).
//   - Validate (func(value []string) error): Custom validation function (default: nil).
//   - Render (func(p *MultiSelectPathPrompt) string): Custom render function (default: nil).
//
// Returns:
//   - *MultiSelectPathPrompt: A new instance of MultiSelectPathPrompt.
func NewMultiSelectPathPrompt(params MultiSelectPathPromptParams) *MultiSelectPathPrompt {
	v := validator.NewValidator("MultiSelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p MultiSelectPathPrompt
	p = MultiSelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[[]string]{
			Context:      params.Context,
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  1,
			Validate:     WrapValidate(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
		Filter:      params.Filter,
		Required:    params.Required,
		FileSystem:  params.FileSystem,
	}

	if cwd, err := p.FileSystem.Getwd(); err == nil && params.InitialPath == "" {
		params.InitialPath = cwd
	}
	p.Root = NewPathNode(params.InitialPath, PathNodeOptions{
		OnlyShowDir: p.OnlyShowDir,
		FileSystem:  p.FileSystem,
	})
	p.CurrentOption = p.Root.FirstChild()
	p.mapSelectedOptions(p.Root)

	actionHandler := NewActionHandler(map[Action]func(){
		UpAction:    func() { p.moveCursor(-1) },
		DownAction:  func() { p.moveCursor(1) },
		LeftAction:  p.closeNode,
		RightAction: p.openNode,
		HomeAction: func() {
			if layerOptions := p.CurrentOption.FilteredLayer(p.Search); len(layerOptions) > 0 {
				p.CurrentOption = layerOptions[0]
				p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
			}
		},
		EndAction: func() {
			if layerOptions := p.CurrentOption.FilteredLayer(p.Search); len(layerOptions) > 0 {
				p.CurrentOption = layerOptions[len(layerOptions)-1]
				p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
			}
		},
		SpaceAction: p.toggleOption,
	}, p.filterOptions)
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))
	})

	p.On(FinalizeEvent, func(args ...any) {
		sort.SliceStable(p.Value, func(i, j int) bool {
			return p.Value[i] < p.Value[j]
		})
	})

	return &p
}

// Options returns a list of filtered PathNode options based on the current layer and search term.
//
// Returns:
//   - []*PathNode: A slice of PathNode objects representing the available options.
func (p *MultiSelectPathPrompt) Options() []*PathNode {
	return p.Root.FilteredFlat(p.Search, p.CurrentOption)
}

// moveCursor moves the cursor up or down within the current layer of options.
//
// Parameters:
//   - direction (int): The direction to move the cursor (-1 for up, 1 for down).
func (p *MultiSelectPathPrompt) moveCursor(direction int) {
	if layerOptions := p.CurrentOption.FilteredLayer(p.Search); len(layerOptions) > 0 {
		layerIndex := p.Root.IndexOf(p.CurrentOption, layerOptions)
		p.CurrentOption = layerOptions[utils.MinMaxIndex(layerIndex+direction, len(layerOptions))]
		p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
	}
}

// closeNode closes the currently selected node or moves up to the parent directory.
func (p *MultiSelectPathPrompt) closeNode() {
	p.Search = ""
	if p.CurrentOption.IsOpen && len(p.CurrentOption.Children) == 0 {
		p.CurrentOption.Close()
		return
	}

	if p.CurrentOption.IsRoot() {
		p.Root = NewPathNode(path.Dir(p.Root.Path), PathNodeOptions{
			OnlyShowDir: p.OnlyShowDir,
			FileSystem:  p.FileSystem,
		})
		p.CurrentOption = p.Root
		p.mapSelectedOptions(p.Root)
		return
	}

	if p.CurrentOption.Parent.IsRoot() {
		p.CurrentOption = p.Root
		return
	}

	p.CurrentOption = p.CurrentOption.Parent
	p.CurrentOption.Close()
}

// openNode opens the currently selected node, revealing its children if any exist.
func (p *MultiSelectPathPrompt) openNode() {
	p.Search = ""
	p.CurrentOption.Open()
	if len(p.CurrentOption.Children) > 0 {
		p.mapSelectedOptions(p.CurrentOption)
		p.CurrentOption = p.CurrentOption.FirstChild()
	}
}

// toggleOption toggles the selection state of the currently selected option.
func (p *MultiSelectPathPrompt) toggleOption() {
	if p.CurrentOption.IsSelected {
		p.CurrentOption.IsSelected = false
		value := []string{}
		for _, v := range p.Value {
			if v != p.CurrentOption.Path {
				value = append(value, v)
			}
		}
		p.Value = value
	} else {
		p.CurrentOption.IsSelected = true
		p.Value = append(p.Value, p.CurrentOption.Path)
	}
}

// filterOptions updates the search term based on the provided key input and filters the available options.
//
// Parameters:
//   - key (*Key): The key event that triggered the filtering.
func (p *MultiSelectPathPrompt) filterOptions(key *Key) {
	if !p.Filter {
		return
	}

	p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
	if p.CurrentOption.IsRoot() {
		return
	}

	layerOptions := p.CurrentOption.FilteredLayer(p.Search)
	layerIndex := p.Root.IndexOf(p.CurrentOption, layerOptions)
	options := p.Options()

	if layerIndex == -1 && len(layerOptions) > 0 {
		p.CurrentOption = layerOptions[0]
	}
	p.CursorIndex = p.Root.IndexOf(p.CurrentOption, options)
}

// mapSelectedOptions traverses the provided node and its children, marking them as selected if their paths match any in the prompt's value.
//
// Parameters:
//   - node (*PathNode): The root node to start traversing from.
func (p *MultiSelectPathPrompt) mapSelectedOptions(node *PathNode) {
	node.TraverseNodes(func(node *PathNode) {
		for _, path := range p.Value {
			if path == node.Path {
				node.IsSelected = true
				return
			}
		}
		node.IsSelected = false
	})
}
