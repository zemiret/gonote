package main

type State struct {
	notebooks []*Notebook
	active    string
}

type Store struct {
	s State
	l *Layout
}

var AppStore *Store

func initState(initial State, layout *Layout) {
	AppStore = &Store{
		s: initial,
		l: layout,
	}
}

func (s *Store) Update(updateFn func(s State) State) error {
	AppStore.s = updateFn(AppStore.s)
	return AppStore.l.Draw()
}
