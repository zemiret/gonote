package main

import (
	"fmt"
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
	var bot, topL, topR *gocui.View
	var botV, topLV, topRV *GridItem

	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

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

		if err == gocui.ErrUnknownView {
			topLV = &GridItem{
				V:      topL,
				Right:  topR,
				Bottom: bot,
				draw: func(maxX, maxY int, v *gocui.View) error {
					_, err := fmt.Fprintf(v, "Testooo left")
					return err
				},
			}

			topRV = &GridItem{
				V:      topR,
				Left:   topL,
				Bottom: bot,
				draw: func(maxX, maxY int, v *gocui.View) error {
					_, err := fmt.Fprintf(v, "Testooo main")
					return err
				},
			}

			botV = &GridItem{
				V:   bot,
				Top: topL,
				draw: func(maxX, maxY int, v *gocui.View) error {
					_, err := fmt.Fprintf(v, "Testooo bot")
					return err
				},
			}

			mainLay.GUI = g

			mainLay.AddView(topRV)
			mainLay.AddView(topLV)
			mainLay.AddView(botV)

			_ = mainLay.Draw()
			_ = mainLay.SetActive(topLV.V.Name())
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
