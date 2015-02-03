package ttt

import (
	"encoding/json"
	"math/rand"
	"time"

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

	Title   = "Tic-tac-toe"
	HelpMsg = `
- 1-PERSON GAME: f1
- 2-PERSON GAME: f2
- LEFT: h, ctrl-b, arrow-left
- DOWN: j, ctrl-n, arrow-down
- UP: k, ctrl-p, arrow-up
- RIGHT: l, ctrl-f, arrow-right
- EXIT: q, esc
- ENTER: i, enter, space
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

var AIOverStatuses = []string{
	StatusWin,
	StatusLoss,
	StatusTie,
	StatusOtherLeft,
}

var Corners = []Position{
	Position{0, 0},
	Position{Size - 1, 0},
	Position{0, Size - 1},
	Position{Size - 1, Size - 1},
}

func RandInt(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

func itemInSlice(i string, s []string) bool {
	for _, item := range s {
		if item == i {
			return true
		}
	}
	return false
}

func IsOverStatus(s string) bool {
	return itemInSlice(s, OverStatuses)
}

func IsAIOverStatus(s string) bool {
	return itemInSlice(s, AIOverStatuses)
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

// Cells in its horizontal row
func (g *Grid) HRowNeighbors(p Position) []Position {
	n := []Position{}
	for i := 0; i < Size; i++ {
		if i != p.X {
			n = append(n, Position{i, p.Y})
		}
	}
	return n
}

// Cells in its vertical row
func (g *Grid) VRowNeighbors(p Position) []Position {
	n := []Position{}
	for i := 0; i < Size; i++ {
		if i != p.Y {
			n = append(n, Position{p.X, i})
		}
	}
	return n
}

// Cells in its left diagonal row (if applicable)
func (g *Grid) LDRowNeighbors(p Position) []Position {
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

// Cells in its right diagonal row (if applicable)
func (g *Grid) RDRowNeighbors(p Position) []Position {
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

// Check if the give position has same marks in a row
func (g *Grid) HasSameMarksInRows(p Position, s string) bool {
	g.Set(p, s)
	ns := [][]Position{
		g.HRowNeighbors(p),
		g.VRowNeighbors(p),
		g.LDRowNeighbors(p),
		g.RDRowNeighbors(p),
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

func (g *Grid) IsEmpty() bool {
	for _, l := range g {
		for _, s := range l {
			if s != "" {
				return false
			}
		}
	}
	return true
}

func (g *Grid) GetRandomCorner() Position {
	return Corners[RandInt(4)]
}

func (g *Grid) GetAvailableCells() []Position {
	var pos []Position
	for x, l := range g {
		for y, s := range l {
			if s == "" {
				pos = append(pos, Position{x, y})
			}
		}
	}
	return pos
}

type GameResult struct {
	Score int
	Pos   Position
}

type GameResults []GameResult

func (gs GameResults) Min() GameResult {
	if len(gs) == 0 {
		return GameResult{}
	}
	min := gs[0]
	for _, g := range gs {
		if g.Score < min.Score {
			min = g
		}
	}
	return min
}

func (gs GameResults) Max() GameResult {
	if len(gs) == 0 {
		return GameResult{}
	}
	max := gs[0]
	for _, g := range gs {
		if g.Score > max.Score {
			max = g
		}
	}
	return max
}

type Game struct {
	CurrentPlayer string
	NextPlayer    string
	Grd           Grid
}

// Return score and whether or not the game is over
func (g Game) Judge(player string, pos Position) (int, bool) {
	if g.Grd.HasSameMarksInRows(pos, player) {
		return 1, true
	} else if g.Grd.IsFull() {
		return 0, true
	} else {
		return 0, false
	}
}

func (g *Game) SwitchTurn() {
	p := g.CurrentPlayer
	g.CurrentPlayer = g.NextPlayer
	g.NextPlayer = p
}

// Use minmax to get the best move for a player
func (g Game) GetBestMove(player string) GameResult {
	if g.Grd.IsEmpty() {
		return GameResult{
			Score: 0,
			Pos:   g.Grd.GetRandomCorner(),
		}
	}
	var gs GameResults

	pos := g.Grd.GetAvailableCells()
	for _, p := range pos {
		ng := g
		score, over := ng.Judge(ng.CurrentPlayer, p)
		if over {
			if ng.CurrentPlayer != player {
				score = -score
			}
			gs = append(gs, GameResult{score, p})
		} else {
			ng.Grd.Set(p, ng.CurrentPlayer)
			(&ng).SwitchTurn()
			rd := ng.GetBestMove(player)
			gs = append(gs, GameResult{rd.Score, p})
		}

	}

	if g.CurrentPlayer == player {
		return gs.Max()
	} else {
		return gs.Min()
	}
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
