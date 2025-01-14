package core

import (
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
	Input        *os.File
	Output       *os.File
	InitialValue string
	OnlyShowDir  bool
	Filter       bool
	FileSystem   FileSystem
	Render       func(p *SelectPathPrompt) string
}

func NewSelectPathPrompt(params SelectPathPromptParams) *SelectPathPrompt {
	v := validator.NewValidator("SelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p SelectPathPrompt
	p = SelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
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

func (p *SelectPathPrompt) Options() []*PathNode {
	return p.Root.FilteredFlat(p.Search, p.CurrentOption)
}

func (p *SelectPathPrompt) moveCursor(direction int) {
	if layerOptions := p.CurrentOption.FilteredLayer(p.Search); len(layerOptions) > 0 {
		layerIndex := p.Root.IndexOf(p.CurrentOption, layerOptions)
		p.CurrentOption = layerOptions[utils.MinMaxIndex(layerIndex+direction, len(layerOptions))]
		p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
	}
}

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

func (p *SelectPathPrompt) openNode() {
	p.Search = ""
	p.CurrentOption.Open()
	if len(p.CurrentOption.Children) > 0 {
		p.CurrentOption = p.CurrentOption.FirstChild()
	}
}

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
