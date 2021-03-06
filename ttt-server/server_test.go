package main

import (
	"container/list"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wujiang/tic-tac-toe"
)

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

func tttsTeardown() {
	group := make(Group)
	players := make(map[string]*Player)
	ttts.Players = &players
	ttts.BenchPlayers = &PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	ttts.WithAIPlayers = make(chan *Player, BufferedChanLen)
	ttts.Announce = make(chan *Announcement, BufferedChanLen)
	ttts.Groups = &group

	aiPlayers := make(map[string]*AIPlayer)
	am.AIPlayers = &aiPlayers
}

func TestTTTScreateNewRound(t *testing.T) {
	player1 := &Player{
		ID:   "player-1",
		Name: "Adam",
	}
	player2 := &Player{
		ID:   "player-2",
		Name: "John",
	}
	ttts.createNewRound(player1, player2)
	assert.Equal(t, len(ttts.Announce), 2)
	a1 := <-ttts.Announce
	a2 := <-ttts.Announce
	assert.True(t, (a1.ToPlayer == *player1 && a1.VSPlayer == *player2 &&
		a2.ToPlayer == *player2 && a2.VSPlayer == *player1) ||
		(a1.ToPlayer == *player2 && a1.VSPlayer == *player1 &&
			a2.ToPlayer == *player1 && a2.VSPlayer == *player2))
	tttsTeardown()
}

func TestTTTSProcessJoinAI(t *testing.T) {
	player1 := &Player{
		ID:   "player-1",
		Name: "Adam",
	}
	ttts.ProcessJoin(player1, true)
	assert.Equal(t, len(*am.AIPlayers), 1)
	assert.Equal(t, len(*ttts.Players), 1)
	assert.Equal(t, len(*ttts.Groups), 1)
	tttsTeardown()
}

func TestTTTSProcessJoinSingle(t *testing.T) {
	player1 := &Player{
		ID:   "player-1",
		Name: "Adam",
	}
	ttts.ProcessJoin(player1, false)
	assert.Equal(t, len(*am.AIPlayers), 0)
	assert.Equal(t, len(*ttts.Players), 1)
	assert.Equal(t, len(*ttts.Groups), 0)
	assert.Equal(t, ttts.BenchPlayers.Len(), 1)
	tttsTeardown()
}

func TestTTTSProcessJoin(t *testing.T) {
	player1 := &Player{
		ID:   "player-1",
		Name: "Adam",
	}
	ttts.ProcessJoin(player1, false)
	player2 := &Player{
		ID:   "player-2",
		Name: "John",
	}
	ttts.ProcessJoin(player2, false)
	assert.Equal(t, len(*am.AIPlayers), 0)
	assert.Equal(t, len(*ttts.Players), 2)
	assert.Equal(t, len(*ttts.Groups), 1)
	assert.Equal(t, ttts.BenchPlayers.Len(), 0)
	tttsTeardown()
}

func TestTTTSProcessQuit(t *testing.T) {
	player1 := &Player{
		ID:   "player-1",
		Name: "Adam",
	}
	ttts.ProcessJoin(player1, false)
	player2 := &Player{
		ID:   "player-2",
		Name: "John",
	}
	ttts.ProcessJoin(player2, false)
	ttts.ProcessQuit(player1)
	assert.Equal(t, len(*am.AIPlayers), 0)
	assert.Equal(t, len(*ttts.Players), 1)
	assert.Equal(t, len(*ttts.Groups), 0)
	assert.Equal(t, ttts.BenchPlayers.Len(), 0)
	tttsTeardown()
}
