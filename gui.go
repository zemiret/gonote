package main

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
)

type Direction int

const (
	Left Direction = iota
	Top
	Right
	Bottom
)

var (
	ErrDrawUnimplemented = errors.New("draw function not supplied")
)

// Draws itself
type Drawer interface {
	Draw(g *gocui.Gui) error
}

// Provides view's name to switch to based on the direction
type Switcher interface {
	Name() string
	SwitchName(d Direction) string
}

// Can draw itself and switch to other
type DrawSwitcher interface {
	Drawer
	Switcher
}

// --- GridItem is a DrawSwitcher ---
type GridItem struct {
	draw   func(maxX, maxY int, v *gocui.View) error
	V      *gocui.View // Main view
	Left   *gocui.View
	Top    *gocui.View
	Right  *gocui.View
	Bottom *gocui.View
}

func (v *GridItem) Name() string {
	return v.V.Name()
}

func (v *GridItem) Draw(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v.draw == nil {
		return ErrDrawUnimplemented
	}

	return v.draw(maxX, maxY, v.V)
}

func (v *GridItem) SwitchName(d Direction) string {
	switch {
	case d == Left && v.Left != nil:
		return v.Left.Name()
	case d == Top && v.Top != nil:
		return v.Top.Name()
	case d == Right && v.Right != nil:
		return v.Right.Name()
	case d == Bottom && v.Bottom != nil:
		return v.Bottom.Name()
	default:
		return ""
	}
}

// --- Layout manages DrawSwitcher in a grid
// In case support for nesting layouts is required, it could implement DrawSwitcher to handle it
type layout struct {
	GUI    *gocui.Gui
	Active DrawSwitcher
	Grid   map[string]DrawSwitcher
}

func NewLayout() *layout {
	return &layout{
		Grid: make(map[string]DrawSwitcher),
	}
}

func (lay *layout) Draw() error {
	for _, view := range lay.Grid {
		if err := view.Draw(lay.GUI); err != nil {
			return err
		}
	}
	return nil
}

func (lay *layout) SwitchActive(d Direction) (err error) {
	newActName := lay.Active.SwitchName(d)
	if newActName != "" {
		err = lay.SetActive(newActName)
	}

	err = nil
	return
}

func (lay *layout) SetActive(name string) error {
	_, err := lay.GUI.SetCurrentView(name)
	if err != nil {
		return fmt.Errorf("error setting current view: %v", err)
	}
	lay.Active = lay.Grid[name]

	return nil
}

func (lay *layout) AddView(v DrawSwitcher) {
	lay.Grid[v.Name()] = v
}
