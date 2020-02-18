package main

type State struct {
	notebooks []*Notebook
	activeIdx int
}

type Store struct {
	s State
	l *Layout
}

var AppStore *Store

func initStore(initial State, layout *Layout) {
	AppStore = &Store{
		s: initial,
		l: layout,
	}
}

func (s *Store) Update(updateFn func(s State) State) error {
	AppStore.s = updateFn(AppStore.s)
	return AppStore.l.Draw()
}

// --------- TODO: SOME MOCK STATE -------------------
var mockState = State{
	activeIdx: 0,
	notebooks: []*Notebook{
		{
			Name: "Notebook1",
			Notes: []Note{
				{
					Name:    "Note11",
					Content: "Note11 content. YO!",
				},
				{
					Name:    "Note12",
					Content: "Note12 content. YO!",
				},
			},
		},
		{
			Name: "Notebook2",
			Notes: []Note{
				{
					Name:    "Note21",
					Content: "Note21 content. YO!",
				},
				{
					Name:    "Note22",
					Content: "Note22 content. YO!",
				},
			},
		},
		{
			Name: "Notebook3",
			Notes: []Note{
				{
					Name:    "Note31",
					Content: "Note21 content. YO!",
				},
			},
		},
	},
}
