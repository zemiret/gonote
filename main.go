package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	mainLay := NewLayout()
	var topL, topR *gocui.View

	notebookList := &NotebookListWidget{
		notebooks: []*Notebook{
			{Name: "Notebook1"},
			{Name: "Notebook2"},
			{Name: "Notebook3"},
		},
		active: "Notebook1",
	}

	var isfirst bool

	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		isfirst = false

		if topL, err = g.SetView("leftPane", 0, 0, 30, maxY-1); err != nil {
			log.Println("Creating left pane")

			if err != gocui.ErrUnknownView {
				panic(err)
			}
			topL.Title = "Top left"

			isfirst = true
		}

		if topR, err = g.SetView("mainPane", 31, 0, maxX-1, maxY-1); err != nil {
			log.Println("Creating main pane")

			if err != gocui.ErrUnknownView {
				panic(err)
			}
			topR.Title = "Top right"
			topR.Editable = true
			topR.Wrap = true
			topR.Autoscroll = true

			isfirst = true
		}

		if isfirst {
			log.Println("isfirst, true")

			topLV := &View{
				View:   topL,
				widget: notebookList,
			}

			topRV := &View{
				View: topR,
			}

			topLGI := &GridItem{
				V:     topLV,
				Right: topRV,
			}

			topRGI := &GridItem{
				V:    topRV,
				Left: topLV,
			}

			mainLay.GUI = g

			mainLay.AddGridItem(topRGI)
			mainLay.AddGridItem(topLGI)

			if err = mainLay.SetActive(topLGI.V.Name()); err != nil {
				log.Panicln(err)
			}
		}

		// TODO: Here we have a culprit. Sometimes this fails. Fix it
		if err = mainLay.Draw(); err != nil {
			log.Panicln(err)
		}

		return nil
	})
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

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

	setKeybinding(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func setKeybinding(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
