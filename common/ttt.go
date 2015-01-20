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
