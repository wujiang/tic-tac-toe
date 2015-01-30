package ttt

import (
	"encoding/json"

	"github.com/nsf/termbox-go"
)

const (
	Width  int = 30
	Height int = 12
	Size   int = 3
	XSpan  int = Width / Size
	YSpan  int = Height / Size

	ColDef = termbox.ColorDefault

	Up    string = "up"
	Down  string = "down"
	Left  string = "left"
	Right string = "right"

	SpecialRune rune = ' '
	MyRune      rune = 'X'
	OtherRune   rune = 'O'

	CmdQuit     string = "Quit"
	CmdJoin     string = "Join"
	CmdJoinAI   string = "Join AI"
	CmdMove     string = "Move"
	CmdNewRound string = "New round"

	StatusInit           string = ""
	StatusConnected      string = "Connected to server"
	StatusWin            string = "You win"
	StatusLoss           string = "You loss"
	StatusTie            string = "Tie"
	StatusQuit           string = "You quit"
	StatusWait           string = "Waiting for another player"
	StatusOtherLeft      string = "The other player left"
	StatusMatched        string = "Matched"
	StatusYourTurn       string = "Your turn"
	StatusWaitTurn       string = "Other user's turn"
	StatusLossConnection string = "Loss connection from server"

	HelpMsg = `
- LEFT: h, ctrl-b, arrow-left
- DOWN: j, ctrl-n, arrow-down
- UP: k, ctrl-p, arrow-up
- RIGHT: l, ctrl-f, arrow-right
- EXIT: q, esc
- ENTER: i, enter, space
- 1 PERSON GAME: f1
- 2 PERSON GAME: f2
`
)

var OverStatuses = []string{
	StatusInit,
	StatusConnected,
	StatusWin,
	StatusLoss,
	StatusTie,
	StatusOtherLeft,
	StatusWait,
}

func IsOverStatus(s string) bool {
	for _, st := range OverStatuses {
		if s == st {
			return true
		}
	}
	return false
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Grid [Size][Size]string

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

func (g *Grid) hRowNeighbors(p Position) []Position {
	n := []Position{}
	for i := 0; i < Size; i++ {
		if i != p.X {
			n = append(n, Position{i, p.Y})
		}
	}
	return n
}

func (g *Grid) vRowNeighbors(p Position) []Position {
	n := []Position{}
	for i := 0; i < Size; i++ {
		if i != p.Y {
			n = append(n, Position{p.X, i})
		}
	}
	return n
}

func (g *Grid) ldRowNeighbors(p Position) []Position {
	n := []Position{}
	isDiagCell := false
	for i := 0; i < Size; i++ {
		lp := Position{i, i}
		if lp == p {
			isDiagCell = true
		} else {
			n = append(n, lp)
		}
	}
	if isDiagCell {
		return n
	} else {
		return []Position{}
	}
}

func (g *Grid) rdRowNeighbors(p Position) []Position {
	n := []Position{}
	isDiagCell := false
	for i := 0; i < Size; i++ {
		rp := Position{i, Size - 1 - i}
		if rp == p {
			isDiagCell = true
		} else {
			n = append(n, rp)
		}
	}
	if isDiagCell {
		return n
	} else {
		return []Position{}
	}
}

func (g *Grid) HasSameMarksInRows(p Position, s string) bool {
	g.Set(p, s)
	ns := [][]Position{
		g.hRowNeighbors(p),
		g.vRowNeighbors(p),
		g.ldRowNeighbors(p),
		g.rdRowNeighbors(p),
	}

	for _, l := range ns {
		if len(l) == 0 {
			continue
		}
		isSame := true
		for _, np := range l {
			if g.Get(p) != g.Get(np) {
				isSame = false
				break
			}
		}
		if isSame {
			return true
		}
	}

	return false
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

type PlayerAction struct {
	RoundID    string   `json:"round_id,omitempty"`
	PlayerID   string   `json:"player_id,omitempty"`
	PlayerName string   `json:"player_name,omitempty"`
	Pos        Position `json:"position"`
	Cmd        string   `json:"cmd"`
}

type PlayerStatus struct {
	RoundID     string `json:"round_id,omitempty"`
	PlayerName  string `json:"player_name,omitempty"`
	PlayerID    string `json:"player_id,omitempty"`
	PlayerScore int    `json:"player_score,omitempty"`
	VSID        string `json:"vs_id,omitempty"`
	VSName      string `json:"vs_name,omitempty"`
	VSScore     int    `json:"score,omitempty"`
	Status      string `json:"status"`
	GridSnap    *Grid  `json:"grid_snap"`
}

func (s *PlayerStatus) Repr() string {
	r, err := json.Marshal(s)
	if err == nil {
		return string(r[:])
	} else {
		return ""
	}
}
