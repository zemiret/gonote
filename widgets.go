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
	SwitchToNoteView ActKind = iota
)

// Describes an action that is to be handled by layout or gui, or some other outside service
type ViewAction struct {
	Kind    ActKind
	Payload interface{}
}

type Widget interface {
	Draw(v *View)
	HandleInput(key interface{}, mod gocui.Modifier, g *gocui.Gui, v *gocui.View) (*ViewAction, error)
	AcceptedInputs() []Input
}

// TODO: Could we make a more generic tree widget out of this? (yes, but notebook need to be changed to some generic interface)

type notebookTreeWidget struct {
	activeNotebookIdx int
	activeNoteIdx     []int
	notebooks         []*Notebook
}

func NotebookTreeWidget(notebooks []*Notebook) *notebookTreeWidget {
	return &notebookTreeWidget{
		activeNotebookIdx: 0,
		activeNoteIdx:     make([]int, len(notebooks)),
		notebooks:         notebooks,
	}
}

// Maybe create something like this for widgets???
//func (nl *notebookTreeWidget) selectState() {
//
//}

func (nl *notebookTreeWidget) Draw(v *View) {
	var err error

	for i, nb := range nl.notebooks {
		if i == nl.activeNotebookIdx {
			_, err = fmt.Fprintln(v, ColorStr(fmt.Sprintf("* %s", nb.Name), 3))
		} else {
			_, err = fmt.Fprintf(v, "* %s\n", nb.Name)
		}

		for j, note := range nl.notebooks[i].Notes {
			if j == nl.activeNoteIdx[i] {
				_, err = fmt.Fprintln(v, ColorStr(fmt.Sprintf("    * %s", note.Name), 5))
			} else {
				_, err = fmt.Fprintf(v, "    * %s\n", note.Name)
			}
		}

		if err != nil {
			panic(err)
		}
	}
}

func (nl *notebookTreeWidget) HandleInput(key interface{}, mod gocui.Modifier, g *gocui.Gui, v *gocui.View) (*ViewAction, error) {
	switch key {
	case gocui.KeyEnter:
		act := &ViewAction{
			Kind:    SwitchToNoteView,
			Payload: nl.notebooks[nl.activeNotebookIdx].Notes[nl.activeNoteIdx[nl.activeNotebookIdx]],
		}
		return act, nil
	case gocui.KeyArrowUp:
		nl.activeNotebookIdx = (nl.activeNotebookIdx - 1 + len(nl.notebooks)) % len(nl.notebooks)
		return nil, nil
	case gocui.KeyArrowDown:
		nl.activeNotebookIdx = (nl.activeNotebookIdx + 1) % len(nl.notebooks)
		return nil, nil
	default:
		return nil, ErrUnsetKey
	}
}

func (nl *notebookTreeWidget) AcceptedInputs() []Input {
	return []Input{
		{Key: gocui.KeyArrowUp, Mod: gocui.ModNone},
		{Key: gocui.KeyArrowDown, Mod: gocui.ModNone},
		{Key: gocui.KeyEnter, Mod: gocui.ModNone},
	}
}
