package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Input struct {
	Key interface{}
	Mod gocui.Modifier
}

type Widget interface {
	Draw(v *View)
	HandleInput(key interface{}, mod gocui.Modifier, g *gocui.Gui, v *gocui.View) error
	AcceptedInputs() []Input
}

type NotebookListWidget struct{}

// Maybe create something like this for widgets???
//func (nl *NotebookListWidget) selectState() {
//
//}

// TODO: Maybe widget should expose all the inputs it can handle?
// Then view could easily read it and add handlers to that automatically (yes, it should)

func (nl *NotebookListWidget) Draw(v *View) {
	var err error

	notebooks := AppStore.s.notebooks
	active := AppStore.s.active

	for _, n := range notebooks {
		if n.Name == active {
			_, err = fmt.Fprintln(v, ColorStr(fmt.Sprintf("* %s", n.Name), 3))
		} else {
			_, err = fmt.Fprintf(v, "* %s\n", n.Name)
		}

		if err != nil {
			panic(err)
		}
	}
}

func (nl *NotebookListWidget) HandleInput(key interface{}, mod gocui.Modifier, g *gocui.Gui, v *gocui.View) error {
	_ = AppStore.Update(func (s State) State {
		s.active = "Notebook2"
		return s
	})
	return nil
}

func (nl *NotebookListWidget) AcceptedInputs() []Input {
	return []Input{
		{Key: gocui.KeyCtrlY, Mod: gocui.ModNone},
	}
}


