package core

import (
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

func NewMultiSelectPathPrompt(params MultiSelectPathPromptParams) *MultiSelectPathPrompt {
	v := validator.NewValidator("MultiSelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p MultiSelectPathPrompt
	p = MultiSelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[[]string]{
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

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})
	p.On(FinalizeEvent, func(args ...any) {
		sort.SliceStable(p.Value, func(i, j int) bool {
			return p.Value[i] < p.Value[j]
		})
	})

	return &p
}

func (p *MultiSelectPathPrompt) Options() []*PathNode {
	return p.Root.FilteredFlat(p.Search, p.CurrentOption)
}

func (p *MultiSelectPathPrompt) handleKeyPress(key *Key) {
	moveCursor := func(direction int) {
		if layerOptions := p.CurrentOption.FilteredLayer(p.Search); len(layerOptions) > 0 {
			layerIndex := p.Root.IndexOf(p.CurrentOption, layerOptions)
			p.CurrentOption = layerOptions[utils.MinMaxIndex(layerIndex+direction, len(layerOptions))]
			p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
		}
	}

	HandleKeyAction(key, map[Action]func(){
		UpAction:   func() { moveCursor(-1) },
		DownAction: func() { moveCursor(1) },
		LeftAction: func() {
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
		},
		RightAction: func() {
			p.Search = ""
			p.CurrentOption.Open()
			if len(p.CurrentOption.Children) == 0 {
				return
			}
			p.mapSelectedOptions(p.CurrentOption)
			p.CurrentOption = p.CurrentOption.FirstChild()
		},
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
		SpaceAction: func() {
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
		},
		DefaultAction: func() {
			if p.Filter {
				p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
				if !p.CurrentOption.IsRoot() {
					layerOptions := p.CurrentOption.FilteredLayer(p.Search)
					layerIndex := p.Root.IndexOf(p.CurrentOption, layerOptions)
					options := p.Options()

					if layerIndex == -1 && len(layerOptions) > 0 {
						p.CurrentOption = layerOptions[0]
					}
					p.CursorIndex = p.Root.IndexOf(p.CurrentOption, options)
				}
			}
		},
	})
}

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
