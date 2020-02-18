package main

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"sync"
)

type Direction int

const (
	Left Direction = iota
	Top
	Right
	Bottom
)

var (
	ErrDrawUnimplemented = errors.New("view does not hav an active widget")
	ErrDrawFail          = errors.New("widget draw fail")
)

func ColorStr(s string, color int) string {
	return fmt.Sprintf("\033[3%d;1m%s\033[0m", color, s)
}

type View struct {
	*gocui.View
	gui *gocui.Gui
	widget Widget
}

func NewView(gui *gocui.Gui, gv *gocui.View, w Widget) (*View, error) {
	v := &View{}

	v.View = gv
	v.gui = gui
	if w != nil {
		return v, v.setWidget(w)
	}
	return v, nil
}

func (v *View) setWidget(w Widget) error {
	v.widget = w
	for _, in := range w.AcceptedInputs() {
		if Config.debug {
			log.Printf("Setting keybinding at %s\n", v.Name())
		}
		if err := v.gui.SetKeybinding(v.Name(), in.Key, in.Mod, func(gui *gocui.Gui, gv *gocui.View) error {
			if Config.debug {
				log.Printf("Binding called for %s\n", v.Name())
			}
			if v.widget != nil {
				return v.widget.HandleInput(in.Key, in.Mod, gui, gv)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

//func (v *View) setKeybinding(g *gocui.Gui, key gocui.Key, mod gocui.Modifier) error {
//	return g.SetKeybinding(v.Name(), key, mod, func(gui *gocui.Gui, gv *gocui.View) (err error) {
//		err = nil
//		if v.widget != nil {
//			err = v.widget.HandleInput(key, mod, gui, gv)
//		}
//		return
//	})
//}

//func widgetInputHandler()

type GridItem struct {
	V      *View // Main view
	Left   *View
	Top    *View
	Right  *View
	Bottom *View
}

func (gi *GridItem) Name() string {
	return gi.V.Name()
}

func (gi *GridItem) Draw(g *gocui.Gui) (derr error) {
	defer func() {
		if err := recover(); err != nil {
			derr = ErrDrawFail
		}
	}()

	if gi.V.widget == nil {
		return ErrDrawUnimplemented
	}

	gi.V.Clear()
	gi.V.widget.Draw(gi.V)
	return
}

func (gi *GridItem) SwitchName(d Direction) string {
	switch {
	case d == Left && gi.Left != nil:
		return gi.Left.Name()
	case d == Top && gi.Top != nil:
		return gi.Top.Name()
	case d == Right && gi.Right != nil:
		return gi.Right.Name()
	case d == Bottom && gi.Bottom != nil:
		return gi.Bottom.Name()
	default:
		return ""
	}
}

// In case support for nesting layouts is required, it could implement DrawSwitcher to handle it
type Layout struct {
	GUI     *gocui.Gui
	Active  *GridItem
	GridMap *GridItemsMap
}

// TODO: Make sure this map really needs concurrent access
type GridItemsMap struct {
	Grid map[string]*GridItem
	sync.RWMutex
}

func NewLayout() *Layout {
	return &Layout{
		GridMap: &GridItemsMap{
			Grid: make(map[string]*GridItem),
		},
	}
}

func (lay *Layout) Draw() error {
	lay.GridMap.RLock()
	defer lay.GridMap.RUnlock()
	for _, view := range lay.GridMap.Grid {
		if Config.debug {
			log.Println("Drawing view ", view.Name())
		}

		if err := view.Draw(lay.GUI); err != nil && err != ErrDrawUnimplemented {
			return err
		}
	}
	return nil
}

func (lay *Layout) SwitchActive(d Direction) (err error) {
	newActName := lay.Active.SwitchName(d)
	if newActName != "" {
		err = lay.SetActive(newActName)
	}

	err = nil
	return
}

func (lay *Layout) SetActive(name string) error {
	_, err := lay.GUI.SetCurrentView(name)
	if err != nil {
		return fmt.Errorf("error setting current view: %v", err)
	}
	lay.GridMap.RLock()
	defer lay.GridMap.RUnlock()
	lay.Active = lay.GridMap.Grid[name]

	return nil
}

func (lay *Layout) AddGridItem(gi *GridItem) {
	if Config.debug {
		log.Println("Adding view: ", gi.Name())
	}

	lay.GridMap.Lock()
	defer lay.GridMap.Unlock()
	lay.GridMap.Grid[gi.Name()] = gi

	if Config.debug {
		log.Println("Gird content---")
		for _, view := range lay.GridMap.Grid {
			log.Println(view.Name())
		}
		log.Println("---")
	}
}
