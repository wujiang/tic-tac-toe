package ttt

import "github.com/nsf/termbox-go"

type Position struct {
	X int
	Y int
}

type Grid [SIZE][SIZE]string

func (g *Grid) Get(p Position) string {
	return g[p.X][p.Y]
}

func (g *Grid) Set(p Position, s string) {
	g[p.X][p.Y] = s
}

// Neighbor cells are defined as:
// - cells in its horizontal row
// - cells in its vertical row
// - cells in its diagonal row (if applicable)
func (g *Grid) NeighborCells(p Position) []Position {
	n := []Position{}
	// horizontal row
	for i := 0; i < SIZE; i++ {
		if i != p.X {
			n = append(n, Position{i, p.Y})
		}
	}
	// vertical row
	for i := 0; i < SIZE; i++ {
		if i != p.Y {
			n = append(n, Position{p.X, i})
		}
	}
	// diagonal row
	diag := []Position{}
	isDiagCell := false
	for i := 0; i < SIZE; i++ {
		p1 := Position{i, i}
		p2 := Position{i, SIZE - 1 - i}
		if p1 == p || p2 == p {
			isDiagCell = true
		}
		if p1 != p {
			diag = append(diag, p1)
		}
		if p2 != p {
			diag = append(diag, p2)
		}
	}
	if isDiagCell {
		n = append(n, diag...)
	}
	return n
}

func (g *Grid) HasSameMarksInRows(p Position, s string) bool {
	g.Set(p, s)
	ns := g.NeighborCells(p)
	same := true
	for _, np := range ns {
		if p != np {
			same = false
			break
		}
	}
	return same
}

const (
	WIDTH  int = 30
	HEIGHT int = 12
	SIZE   int = 3
	XSPAN  int = WIDTH / SIZE
	YSPAN  int = HEIGHT / SIZE

	COLDEF = termbox.ColorDefault

	UP    string = "up"
	DOWN  string = "down"
	LEFT  string = "left"
	RIGHT string = "right"

	SPECIALRUNE rune = ' '
	MYRUNE      rune = 'O'
	OTHERRUNE   rune = 'X'

	HELPMSG = `
Tic-tac-toe manual:
- LEFT: h, ctrl-b, arrow-left
- DOWN: j, ctrl-n, arrow-down
- UP: k, ctrl-p, arrow-up
- RIGHT: l, ctrl-f, arrow-right
- EXIT: q, esc
- ENTER: o, enter
`
)
