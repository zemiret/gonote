package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type Direction int

const (
	Left Direction = iota
	Top
	Right
	Bottom
)

type View struct {
	V      *gocui.View
	Left   *gocui.View
	Top    *gocui.View
	Right  *gocui.View
	Bottom *gocui.View
}

func (v *View) FocusName(d Direction) string {
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

type Layout struct {
	G      *gocui.Gui
	Active *View
	Grid   map[string]*View
}

func NewLayout() *Layout {
	return &Layout{
		Grid: make(map[string]*View),
	}
}

func (lay *Layout) SwitchActive(d Direction) (err error) {
	newActName := lay.Active.FocusName(d)
	if newActName != "" {
		err = lay.SetActive(newActName)
	}

	err = nil
	return
}

func (lay *Layout) SetActive(name string) error {
	_, err := lay.G.SetCurrentView(name)
	if err != nil {
		return fmt.Errorf("error setting current view: %v", err)
	}
	lay.Active = lay.Grid[name]

	return nil
}

func (lay *Layout) AddView(v *View) {
	lay.Grid[v.V.Name()] = v
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//g.SetManagerFunc(layout)

	// TODO: IMPORTANT! While the below code works (adding view outside of LayoutManager)
	// They are added, well, outside of the manager.
	// Thus, e.g. windows resize is not being handled, because they don't seem to be registered in the manager.
	// Conclusion: Figure out the MY LAYOUT creation process maintaining gocui LayoutManager
	g.SetManagerFunc(func(_ *gocui.Gui) error { return nil })
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	maxX, maxY := g.Size()

	var bot, topL, topR *gocui.View

	if topL, err = g.SetView("leftPane", 0, 0, 30, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
		topL.Title = "Top left"
	}

	if topR, err = g.SetView("mainPane", 31, 0, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
		topR.Title = "Top right"
		topR.Editable = true
		topR.Wrap = true
		topR.Autoscroll = true
	}

	if bot, err = g.SetView("bottomPane", 0, maxY-7, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
		bot.Title = "Bottom pane"
	}

	mainLay := NewLayout()
	mainLay.G = g

	botV := &View{
		V:   bot,
		Top: topL,
	}

	topLV := &View{
		V:      topL,
		Right:  topR,
		Bottom: bot,
	}

	topRV := &View{
		V:      topR,
		Left:   topL,
		Bottom: bot,
	}

	mainLay.AddView(botV)
	mainLay.AddView(topLV)
	mainLay.AddView(topRV)
	_ = mainLay.SetActive(botV.V.Name())

	setKeybinding(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		_ = mainLay.SwitchActive(Left)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		_ = mainLay.SwitchActive(Bottom)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		_ = mainLay.SwitchActive(Top)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		_ = mainLay.SwitchActive(Right)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("leftPane", 0, 0, 30, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Notes"
	}

	if v, err := g.SetView("mainPane", 31, 0, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Main"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("bottomPane", 0, maxY-7, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Bottom pane"
	}

	if _, err := g.SetCurrentView("mainPane"); err != nil {
		return err
	}

	return nil
}

//func layout2(g *gocui.Gui) error {
//	if v, err := g.SetView("top", 20, 15, 47, 25); err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//		v.Title = "TOP"
//	}
//
//	return nil
//}

func setKeybinding(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
