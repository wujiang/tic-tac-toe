package ttt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandInt(t *testing.T) {
	n := RandInt(100)
	assert.True(t, n < 100 && n >= 0)
}

func TestitemInSlice(t *testing.T) {
	assert.True(t, itemInSlice("hello", []string{"hello", "world"}))
	assert.False(t, itemInSlice("foo", []string{"hello", "world"}))
}

func TestGridGet(t *testing.T) {
	gd := Grid{}
	assert.Equal(t, gd.Get(Position{0, 0}), "")
}

func TestGridSet(t *testing.T) {
	gd := Grid{}
	pos := Position{0, 0}
	gd.Set(pos, "X")
	assert.Equal(t, gd.Get(pos), "X")
}

func TestGridHRowNeighbors(t *testing.T) {
	gd := Grid{}
	neighbors := gd.HRowNeighbors(Position{1, 1})
	assert.Equal(t, neighbors, []Position{Position{0, 1}, Position{2, 1}})
}

func TestGridVRowNeighbors(t *testing.T) {
	gd := Grid{}
	neighbors := gd.VRowNeighbors(Position{1, 1})
	assert.Equal(t, neighbors, []Position{Position{1, 0}, Position{1, 2}})
}

func TestGridLDRowNeighbors(t *testing.T) {
	gd := Grid{}
	neighbors := gd.LDRowNeighbors(Position{1, 1})
	assert.Equal(t, neighbors, []Position{Position{0, 0}, Position{2, 2}})
	neighbors = gd.LDRowNeighbors(Position{0, 1})
	assert.Equal(t, neighbors, []Position{})
}

func TestGridRDRowNeighbors(t *testing.T) {
	gd := Grid{}
	neighbors := gd.RDRowNeighbors(Position{1, 1})
	assert.Equal(t, neighbors, []Position{Position{0, 2}, Position{2, 0}})
	neighbors = gd.RDRowNeighbors(Position{0, 1})
	assert.Equal(t, neighbors, []Position{})
}

func TestGridHasSameMarksInRows(t *testing.T) {
	gd := Grid{}
	assert.True(t, gd.HasSameMarksInRows(Position{1, 1}, ""))
	assert.False(t, gd.HasSameMarksInRows(Position{1, 1}, "X"))
	gd = Grid{
		{"X", "O", "X"},
		{"", "", "O"},
		{"O", "", "X"},
	}
	assert.True(t, gd.HasSameMarksInRows(Position{1, 1}, "X"))
	assert.False(t, gd.HasSameMarksInRows(Position{1, 0}, "X"))
}

func TestGridIsFull(t *testing.T) {
	gd := Grid{}
	assert.False(t, gd.IsFull())
	gd = Grid{
		{"X", "O", "X"},
		{"", "", "O"},
		{"O", "", "X"},
	}
	assert.False(t, gd.IsFull())
	gd = Grid{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", "X"},
	}
	assert.True(t, gd.IsFull())
}

func TestGridIsEmpty(t *testing.T) {
	gd := Grid{}
	assert.True(t, gd.IsEmpty())
	gd = Grid{
		{"X", "O", "X"},
		{"", "", "O"},
		{"O", "", "X"},
	}
	assert.False(t, gd.IsEmpty())
}

func TestGridGetAvailableCells(t *testing.T) {
	gd := Grid{
		{"X", "", "O"},
		{"", "O", "X"},
		{"X", "", "O"},
	}
	available := []Position{Position{0, 1}, Position{1, 0}, Position{2, 1}}
	assert.Equal(t, gd.GetAvailableCells(), available)
}

func TestGameResultsMin(t *testing.T) {
	gs := GameResults{}
	assert.Equal(t, gs.Min(), GameResult{})

	gs = GameResults{
		GameResult{1, Position{1, 2}},
		GameResult{9, Position{2, 2}},
		GameResult{-8, Position{0, 2}},
		GameResult{10, Position{0, 1}},
	}
	assert.Equal(t, gs.Min(), gs[2])
}

func TestGameResultsMax(t *testing.T) {
	gs := GameResults{}
	assert.Equal(t, gs.Max(), GameResult{})

	gs = GameResults{
		GameResult{1, Position{1, 2}},
		GameResult{9, Position{2, 2}},
		GameResult{-8, Position{0, 2}},
		GameResult{10, Position{0, 1}},
	}
	assert.Equal(t, gs.Max(), gs[3])
}

func TestGameSwitchTurn(t *testing.T) {
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           Grid{},
	}
	g.SwitchTurn()
	assert.Equal(t, g.CurrentPlayer, "O")
	assert.Equal(t, g.NextPlayer, "X")
}

func TestGameJudgeNotOver(t *testing.T) {
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           Grid{},
	}
	score, over := g.Judge("X", Position{0, 0})
	assert.Equal(t, score, 0)
	assert.False(t, over)
}

func TestGameJudgeWin(t *testing.T) {
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           Grid{{"X", "X"}},
	}
	score, over := g.Judge("X", Position{0, 2})
	assert.Equal(t, score, 1)
	assert.True(t, over)
}

func TestGameJudgeTie(t *testing.T) {
	grid := Grid{
		{"X", "X", "O"},
		{"O", "O", "X"},
		{"", "X", "O"},
	}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	score, over := g.Judge("X", Position{2, 0})
	assert.Equal(t, score, 0)
	assert.True(t, over)
}

func TestGameGetBestMove(t *testing.T) {
	grid := Grid{
		{"O", "X", "X"},
		{"", "", "O"},
		{"X", "", "O"},
	}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	assert.Equal(t, r, GameResult{1, Position{1, 1}})
}

func TestGameGetBestMove2(t *testing.T) {
	grid := Grid{}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	assert.Equal(t, r.Score, 0)
	assert.True(t, Corners[0] == r.Pos || Corners[1] == r.Pos ||
		Corners[2] == r.Pos || Corners[3] == r.Pos)
}

func TestGameGetBestMove3(t *testing.T) {
	grid := Grid{
		{"X", "", ""},
		{"", "", ""},
		{"", "", ""},
	}
	g := Game{
		CurrentPlayer: "O",
		NextPlayer:    "X",
		Grd:           grid,
	}
	r := g.GetBestMove("O")
	assert.Equal(t, r, GameResult{0, Position{1, 1}})
}

func TestGameGetBestMove4(t *testing.T) {
	grid := Grid{
		{"X", "", ""},
		{"", "O", ""},
		{"", "", ""},
	}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	assert.Equal(t, r, GameResult{0, Position{0, 1}})
}

func TestGameGetBestMove5(t *testing.T) {
	grid := Grid{
		{"X", "X", ""},
		{"", "O", ""},
		{"", "", ""},
	}
	g := Game{
		CurrentPlayer: "O",
		NextPlayer:    "X",
		Grd:           grid,
	}
	r := g.GetBestMove("O")
	assert.Equal(t, r, GameResult{0, Position{0, 2}})
}

func TestGameGetBestMove6(t *testing.T) {
	grid := Grid{
		{"X", "X", "O"},
		{"", "O", ""},
		{"", "", ""},
	}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	assert.Equal(t, r, GameResult{0, Position{2, 0}})
}

func TestGameGetBestMove7(t *testing.T) {
	grid := Grid{
		{"X", "X", "O"},
		{"", "O", ""},
		{"X", "", ""},
	}
	g := Game{
		CurrentPlayer: "O",
		NextPlayer:    "X",
		Grd:           grid,
	}
	r := g.GetBestMove("O")
	assert.Equal(t, r, GameResult{0, Position{1, 0}})
}

func TestGameGetBestMove8(t *testing.T) {
	grid := Grid{
		{"X", "X", "O"},
		{"O", "O", ""},
		{"X", "", ""},
	}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	assert.Equal(t, r, GameResult{0, Position{1, 2}})
}

func TestGameGetBestMove9(t *testing.T) {
	grid := Grid{
		{"X", "X", "O"},
		{"O", "O", "X"},
		{"X", "", ""},
	}
	g := Game{
		CurrentPlayer: "O",
		NextPlayer:    "X",
		Grd:           grid,
	}
	r := g.GetBestMove("O")
	assert.Equal(t, r, GameResult{0, Position{2, 1}})
}

func TestGameGetBestMove10(t *testing.T) {
	grid := Grid{
		{"X", "X", "O"},
		{"O", "O", "X"},
		{"X", "O", ""},
	}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	assert.Equal(t, r, GameResult{0, Position{2, 2}})
}

func TestGameGetBestMove11(t *testing.T) {
	grid := Grid{
		{"O", "X", "X"},
		{"X", "", "O"},
		{"X", "", "O"},
	}
	g := Game{
		CurrentPlayer: "O",
		NextPlayer:    "X",
		Grd:           grid,
	}
	r := g.GetBestMove("O")
	assert.Equal(t, r, GameResult{1, Position{1, 1}})
}
