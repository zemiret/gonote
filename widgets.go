package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Widget interface {
	Draw(v *View)
}

type NotebookListWidget struct {
	notebooks []*Notebook
	active    string
}

func (nl *NotebookListWidget) Draw(v *View) {
	for _, n := range nl.notebooks {
		if n.Name == nl.active {
			v.FgColor = gocui.ColorGreen
		}
		v.FgColor = gocui.ColorDefault

		_, err := fmt.Fprintf(v, "* %s\n", n.Name)
		if err != nil {
			panic(err)
		}
	}
}


