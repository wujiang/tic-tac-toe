package main

import (
	"container/list"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wujiang/tic-tac-toe"
)

var server *httptest.Server

func setup() {
	// server = httptest.NewServer(WSHandler)
}

func teardown() {
	server.Close()
}

func TestPlayersQueuePush(t *testing.T) {
	pq := PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	player := &Player{}
	assert.Equal(t, pq.players.Len(), 0)
	pq.Push(player)
	assert.Equal(t, pq.players.Len(), 1)
	pq.Push(player)
	assert.Equal(t, pq.players.Len(), 2)
}

func TestPlayersQueuePop(t *testing.T) {
	pq := PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	player := &Player{}
	pq.players.PushBack(player)
	assert.Equal(t, pq.players.Len(), 1)
	assert.Equal(t, pq.Pop(), player)
	assert.Equal(t, pq.players.Len(), 0)
}

func TestPlayersQueueLen(t *testing.T) {
	pq := PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	player := &Player{}
	pq.Push(player)
	assert.Equal(t, pq.Len(), 1)
	assert.Equal(t, pq.Pop(), player)
	assert.Equal(t, pq.Len(), 0)
}

func TestPlayersQueueRemove(t *testing.T) {
	pq := PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	player := &Player{}
	pq.Push(player)
	assert.Equal(t, pq.Len(), 1)
	pq.Remove(player)
	assert.Equal(t, pq.Len(), 0)
}

func TestRoundswitchTurn(t *testing.T) {
	player1 := &Player{
		ID: "Adam",
	}
	player2 := &Player{
		ID: "John",
	}
	rd := Round{
		CurrentPlayer: player1,
		NextPlayer:    player2,
	}
	rd.switchTurn()
	assert.True(t, rd.CurrentPlayer == player2 && rd.NextPlayer == player1)
}

func TestRoundgetOtherPlayer(t *testing.T) {
	player1 := &Player{
		ID: "Adam",
	}
	player2 := &Player{
		ID: "John",
	}
	rd := &Round{
		CurrentPlayer: player1,
		NextPlayer:    player2,
	}
	assert.Nil(t, rd.getOtherPlayer(&Player{}))
	assert.Equal(t, rd.getOtherPlayer(player2), player1)
	assert.Equal(t, rd.getOtherPlayer(player1), player2)
}

func TestAnnouncementtoPlayerStatus(t *testing.T) {
	player1 := Player{
		ID: "Adam",
	}
	round := Round{
		ID:            "round-id",
		CurrentPlayer: &player1,
	}
	ann := &Announcement{
		ToPlayer: player1,
		Rd:       round,
		Status:   ttt.StatusWait,
	}

	expected := ttt.PlayerStatus{
		RoundID:     ann.Rd.ID,
		PlayerID:    ann.ToPlayer.ID,
		PlayerScore: ann.ToPlayer.Score,
		Status:      ann.Status,
		GridSnap:    ann.Rd.Grid,
	}
	ps := ann.toPlayerStatus()
	assert.Equal(t, *ps, expected)

	player2 := Player{
		ID: "John",
	}
	ann.VSPlayer = player2
	expected = ttt.PlayerStatus{
		RoundID:     ann.Rd.ID,
		PlayerID:    ann.ToPlayer.ID,
		PlayerScore: ann.ToPlayer.Score,
		VSID:        ann.VSPlayer.ID,
		VSName:      ann.VSPlayer.Name,
		VSScore:     ann.VSPlayer.Score,
		Status:      ann.Status,
		GridSnap:    ann.Rd.Grid,
	}
	ps = ann.toPlayerStatus()
	assert.Equal(t, *ps, expected)
}
