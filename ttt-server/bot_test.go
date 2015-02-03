package main

import (
	"container/list"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wujiang/tic-tac-toe"
)

func amSetup() {
	am.NewAIPlayer("bot1")
}

func amTeardown() {
	group := make(Group)
	players := make(map[string]*Player)
	ttts.Players = &players
	ttts.BenchPlayers = &PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	ttts.AIPlayers = make(chan *Player, BufferedChanLen)
	ttts.Announce = make(chan *Announcement, BufferedChanLen)
	ttts.Groups = &group

	aiPlayers := make(map[string]*AIPlayer)
	am.AIPlayers = &aiPlayers
}

func TestAIPlayerUpdate(t *testing.T) {
	amSetup()
	ps := &ttt.PlayerStatus{
		RoundID:     "round-1",
		PlayerName:  "AI",
		PlayerID:    "bot1",
		PlayerScore: 1,
		VSID:        "player-1",
		VSName:      "Adam",
		VSScore:     -1,
		Status:      ttt.StatusWait,
	}
	ap := (*am.AIPlayers)["bot1"]
	ap.Update(ps)
	assert.Equal(t, ap.Score, 1)
	assert.Equal(t, ap.VSID, "player-1")
	assert.Equal(t, ap.VSName, "Adam")
	assert.Equal(t, ap.VSScore, -1)
	assert.Equal(t, ap.Status, ttt.StatusWait)
	amTeardown()
}
