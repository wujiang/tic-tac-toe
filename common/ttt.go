package ttt

import "github.com/nsf/termbox-go"

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

	CMD_QUIT string = "QUIT"
	CMD_JOIN string = "JOIN"
	CMD_MOVE string = "MOVE"

	STATUS_WIN       string = "You win"
	STATUS_LOSS      string = "You loss"
	STATUS_TIE       string = "Tie"
	STATUS_QUIT      string = "You quit"
	STATUS_WAIT      string = "Waiting for another player"
	STATUS_LEFT      string = "The other play left"
	STATUS_MATCHED   string = "Matched"
	STATUS_YOUR_TURN string = "Your turn"
	STATUS_WAIT_TURN string = "Other user's turn"

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

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Grid [SIZE][SIZE]string

type PlayerAction struct {
	RoundID    string   `json:"round_id,omitempty"`
	PlayerName string   `json:"player_name,omitempty"`
	Pos        Position `json:"position"`
	Cmd        string   `json:"cmd"`
}

type PlayerStatus struct {
	RoundID    string `json:"round_id,omitempty"`
	PlayerName string `json:"player_name,omitempty"`
	VSName     string `json:"vs_name,omitempty"`
	Status     string `json:"status"`
	GridSnap   *Grid  `json:"grid_snap"`
}

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

func (g *Grid) IsFull() bool {
	for _, l := range g {
		for _, s := range l {
			if s == "" {
				return false
			}
		}
	}
	return true
}
