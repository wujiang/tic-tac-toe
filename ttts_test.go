package ttt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundsMin(t *testing.T) {
	rs := Rounds{}
	assert.Equal(t, rs.Min(), Round{})

	rs = Rounds{
		Round{1, Position{1, 2}},
		Round{9, Position{2, 2}},
		Round{-8, Position{0, 2}},
		Round{10, Position{0, 1}},
	}
	assert.Equal(t, rs.Min(), rs[2])
}

func TestRoundsMax(t *testing.T) {
	rs := Rounds{}
	assert.Equal(t, rs.Max(), Round{})

	rs = Rounds{
		Round{1, Position{1, 2}},
		Round{9, Position{2, 2}},
		Round{-8, Position{0, 2}},
		Round{10, Position{0, 1}},
	}
	assert.Equal(t, rs.Max(), rs[3])
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
	assert.Equal(t, r, Round{1, Position{1, 1}})
}

func TestGameGetBestMove2(t *testing.T) {
	grid := Grid{}
	g := Game{
		CurrentPlayer: "X",
		NextPlayer:    "O",
		Grd:           grid,
	}
	r := g.GetBestMove("X")
	corners := []Position{
		Position{0, 0},
		Position{Size - 1, 0},
		Position{0, Size - 1},
		Position{Size - 1, Size - 1},
	}

	assert.Equal(t, r.Score, 0)
	assert.True(t, corners[0] == r.Pos || corners[1] == r.Pos ||
		corners[2] == r.Pos || corners[3] == r.Pos)
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
	assert.Equal(t, r, Round{0, Position{1, 1}})
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
	assert.Equal(t, r, Round{0, Position{0, 1}})
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
	assert.Equal(t, r, Round{0, Position{0, 2}})
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
	assert.Equal(t, r, Round{0, Position{2, 0}})
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
	assert.Equal(t, r, Round{0, Position{1, 0}})
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
	assert.Equal(t, r, Round{0, Position{1, 2}})
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
	assert.Equal(t, r, Round{0, Position{2, 1}})
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
	assert.Equal(t, r, Round{0, Position{2, 2}})
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
	assert.Equal(t, r, Round{1, Position{1, 1}})
}
