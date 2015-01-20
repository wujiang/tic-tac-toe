package server

import (
	"container/list"
	"net/http"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe/common"
)

type Conn struct {
	ws   *websocket.Conn
	send chan []PlayerMove
}

type Player struct {
	Name string
	Cn   Conn
}

type Round struct {
	RoundID  string
	P1       *Player
	P2       *Player
	NextTurn *Player
	Winner   *Player
	Grid     *ttt.Grid
}

type PlayerMove struct {
	RoundID string
	Plyer   *Player
	Pos     ttt.Position
}

type TTTServer struct {
	Groups       *map[string]Round
	BenchPlayers *list.List

	Joins        chan []Player     // incoming channel
	QuitWaitings chan []Player     // incoming channel
	QuitRounds   chan []PlayerMove // incoming channel
	PlayerMoves  chan []PlayerMove // incoming channel
	Rounds       chan []Round      // outgoing channel
}

func (ttts *TTTServer) MakeRound(p Player) {
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	ttts.BenchPlayers.PushBack(p)
	if ttts.BenchPlayers.Len() > 1 {
		p1 := ttts.BenchPlayers.Front()
		ttts.BenchPlayers.Remove(p1)
		p2 := ttts.BenchPlayers.Front()
		ttts.BenchPlayers.Remove(p2)
		r := Round{}
		r.RoundID = uuid.New()
		r.P1 = &p1
		r.P2 = &p2
		r.NextTurn = &p1
		ttts.Groups[r.RoundID] = r
	}
}

func (ttts *TTTServer) QuitWaiting(p Player) {
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	for e := ttts.BenchPlayers.Front(); e != nil; e = e.Next() {
		if e == p {
			ttts.BenchPlayers.Remove(e)
			break
		}
	}
}

// Remove round from Groups and put the other user into BenchPlayers
func (ttts *TTTServer) QuitRound(m PlayerMove) {
	rd := (*ttts.Groups)[m.RoundID]
	if rd == nil {
		return
	}
	toNotify := rd.P1
	if rd.P1 == m.Plyer {
		toNotify = rd.P2
		rd.P1 = nil
	} else {
		rd.P2 = nil
	}
	ttts.Rounds <- rd

	// TODO: remove round from groups

	ttts.MakeRound(toNotify)
}

// TODO: ttts.Rounds may have a nil player, which means the other player is gone

func (ttts *TTTServer) Judge(m PlayerMove) Round {
	rd := (*ttts.Groups)[m.RoundID]
	if rd.Grid.HasSameMarksInRows(m.Pos, m.Plyer.Name) {
		rd.Winner = m.Plyer
	}
	rd.NextTurn = m.Plyer
	(*ttts.Groups)[m.RoundID] = rd
	return rd
}

func (ttts *TTTServer) run() {
	for {
		select {
		case m := <-ttts.PlayerMoves:
			ttts.Rounds <- ttts.Judge(m)
		case j := <-ttts.Joins:
			ttts.MakeRound(j)
		case w := <-ttts.QuitWaitings:
			ttts.QuitWaiting(w)
		case r := <-ttts.QuitRounds:
			ttts.QuitRound(r)
		}

	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  READ_BUFFER_SIZE,
	WriteBufferSize: WRITE_BUFFER_SIZE,
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
}
