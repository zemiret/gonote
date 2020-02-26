package main

type Note struct {
	Name string
	Content string
}

type Notebook struct {
	Notes []*Note
	Name string
}
