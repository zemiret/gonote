package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Input struct {
	Key interface{}
	Mod gocui.Modifier
}

// TODO: Try to come up with these "higher level" actions to some sensible solution
type ActKind int

const (
	SwitchView ActKind = iota
)

// Describes an action that is to be handled by layout or gui, or some other outside service
type Action struct {
	Kind    ActKind
	Payload interface{}
}

type Widget interface {
	Draw(v *View)
	HandleInput(key interface{}, mod gocui.Modifier, g *gocui.Gui, v *gocui.View) error
	AcceptedInputs() []Input
}

type NotebookListWidget struct{}

// TODO: Could we make a more generic tree widget out of this?

// Maybe create something like this for widgets???
//func (nl *NotebookListWidget) selectState() {
//
//}

func (nl *NotebookListWidget) Draw(v *View) {
	var err error

	notebooks := AppStore.s.notebooks
	active := AppStore.s.activeIdx

	for i, n := range notebooks {
		if i == active {
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
	switch key {
	case gocui.KeyArrowUp:
		return AppStore.Update(func(s State) State {
			s.activeIdx = (s.activeIdx - 1 + len(s.notebooks)) % len(s.notebooks)
			return s
		})
	case gocui.KeyArrowDown:
		return AppStore.Update(func(s State) State {
			s.activeIdx = (s.activeIdx + 1) % len(s.notebooks)
			return s
		})
	default:
		return ErrUnsetKey
	}
}

func (nl *NotebookListWidget) AcceptedInputs() []Input {
	return []Input{
		{Key: gocui.KeyArrowUp, Mod: gocui.ModNone},
		{Key: gocui.KeyArrowDown, Mod: gocui.ModNone},
	}
}
