package core

import (
	"context"
	"os"
	"path"

	"github.com/Mist3rBru/go-clack/core/internals"
	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type SelectPathPrompt struct {
	Prompt[string]
	Root          *PathNode
	CurrentLayer  []*PathNode
	CurrentOption *PathNode
	OnlyShowDir   bool
	Search        string
	Filter        bool
	FileSystem    FileSystem
}

type SelectPathPromptParams struct {
	Context      context.Context
	Input        *os.File
	Output       *os.File
	InitialValue string
	OnlyShowDir  bool
	Filter       bool
	FileSystem   FileSystem
	Render       func(p *SelectPathPrompt) string
}

// NewSelectPathPrompt initializes and returns a new instance of SelectPathPrompt.
//
// The user can navigate through directories and files using arrow keys.
// The user can select an option using enter key.
// The prompt returns the path of the selected option.
// If the user cancels the prompt, it returns an error.
// If an error occurs during the prompt, it also returns an error.
//
// Parameters:
//   - Context (context.Context): The context for the prompt (default: context.Background).
//   - Input (*os.File): The input stream for the prompt (default: OSFileSystem).
//   - Output (*os.File): The output stream for the prompt (default: OSFileSystem).
//   - InitialValue (string): The initial path value (default: current working directory).
//   - OnlyShowDir (bool): Whether to only show directories (default: false).
//   - Filter (bool): Whether to enable filtering of options (default: false).
//   - FileSystem (FileSystem): The file system implementation to use (default: OSFileSystem).
//   - Render (func(p *SelectPathPrompt) string): Custom render function for the prompt (default: nil).
//
// Returns:
//   - *SelectPathPrompt: A new instance of SelectPathPrompt.
func NewSelectPathPrompt(params SelectPathPromptParams) *SelectPathPrompt {
	v := validator.NewValidator("SelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p SelectPathPrompt
	p = SelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Context:     params.Context,
			Input:       params.Input,
			Output:      params.Output,
			CursorIndex: 1,
			Render:      WrapRender[string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
		Filter:      params.Filter,
		FileSystem:  params.FileSystem,
	}

	if cwd, err := p.FileSystem.Getwd(); err == nil && params.InitialValue == "" {
		params.InitialValue = cwd
	}
	p.Root = NewPathNode(params.InitialValue, PathNodeOptions{
		OnlyShowDir: p.OnlyShowDir,
		FileSystem:  p.FileSystem,
	})
	p.CurrentLayer = p.Root.Children
	p.CurrentOption = p.Root.Children[0]
	p.Value = p.CurrentOption.Path

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
	}, p.filterOptions)
	p.On(KeyEvent, func(args ...any) {
		actionHandler(args[0].(*Key))

		if p.CurrentOption != nil {
			p.Value = p.CurrentOption.Path
		} else {
			p.Value = *new(string)
		}
	})

	return &p
}

// Options returns a list of filtered and flattened PathNode options based on the current search term and selected node.
//
// Returns:
//   - []*PathNode: A slice of PathNode objects representing the available options.
func (p *SelectPathPrompt) Options() []*PathNode {
	return p.Root.FilteredFlat(p.Search, p.CurrentOption)
}

// moveCursor moves the cursor up or down within the current layer of options.
//
// Parameters:
//   - direction (int): The direction to move the cursor (-1 for up, 1 for down).
func (p *SelectPathPrompt) moveCursor(direction int) {
	if layerOptions := p.CurrentOption.FilteredLayer(p.Search); len(layerOptions) > 0 {
		layerIndex := p.Root.IndexOf(p.CurrentOption, layerOptions)
		p.CurrentOption = layerOptions[utils.MinMaxIndex(layerIndex+direction, len(layerOptions))]
		p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
	}
}

// closeNode closes the currently selected node or moves up to the parent directory.
func (p *SelectPathPrompt) closeNode() {
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
func (p *SelectPathPrompt) openNode() {
	p.Search = ""
	p.CurrentOption.Open()
	if len(p.CurrentOption.Children) > 0 {
		p.CurrentOption = p.CurrentOption.FirstChild()
	}
}

// filterOptions updates the search term based on the provided key input and filters the available options.
//
// Parameters:
//   - key (*Key): The key event that triggered the filtering.
func (p *SelectPathPrompt) filterOptions(key *Key) {
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
