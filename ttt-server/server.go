package main

import (
	"container/list"
	"net/http"
	"strings"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe"
)

const (
	ReadBufferSize  int = 1024
	WriteBufferSize int = 2048
	BufferedChanLen int = 10
)

var ttts = InitTTTServer()

// Internal communication channel with AI bots
var playerActions = make(chan ttt.PlayerAction, BufferedChanLen)
var playerStatuses = make(chan ttt.PlayerStatus, BufferedChanLen)

type Player struct {
	WS      *websocket.Conn
	RoundID string
	ID      string
	Name    string
	Score   int
}

func (p *Player) repr() string {
	return p.Name + " (" + p.ID + ")"
}

// Parse the action sent by a client
func (p *Player) parseAction() {
	for {
		m := ttt.PlayerAction{}
		p.WS.ReadJSON(&m)
		switch m.Cmd {
		case ttt.CmdQuit:
			ttts.ProcessQuit(p)
			return
		case ttt.CmdJoin:
			p.Name = m.PlayerName
			ttts.ProcessJoin(p, false)
		case ttt.CmdJoinAI:
			p.Name = m.PlayerName
			ttts.ProcessJoin(p, true)
		case ttt.CmdMove:
			ttts.Judge(&m)
		}
	}
}

type PlayersQueue struct {
	players *list.List
	lock    sync.Mutex
}

func (q *PlayersQueue) Pop() *Player {
	q.lock.Lock()
	defer q.lock.Unlock()
	head := q.players.Front()
	if head != nil {
		v := q.players.Remove(head)
		return v.(*Player)
	}
	return nil
}

func (q *PlayersQueue) Push(p *Player) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.players.PushBack(p)
}

func (q *PlayersQueue) Len() int {
	return q.players.Len()
}

func (q *PlayersQueue) Remove(p *Player) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for e := q.players.Front(); e != nil; e = e.Next() {
		if e.Value.(*Player) == p {
			q.players.Remove(e)
			break
		}
	}
}

type Round struct {
	ID            string
	CurrentPlayer *Player
	NextPlayer    *Player
	Winner        *Player
	Grid          *ttt.Grid
}

// Switch turn in a matching round
func (r *Round) switchTurn() {
	temp := r.CurrentPlayer
	r.CurrentPlayer = r.NextPlayer
	r.NextPlayer = temp
}

func (r *Round) getOtherPlayer(p *Player) *Player {
	if r.CurrentPlayer != p && r.NextPlayer == p {
		return r.CurrentPlayer
	} else if r.NextPlayer != p && r.CurrentPlayer == p {
		return r.NextPlayer
	} else {
		return nil
	}
}

type Announcement struct {
	ToPlayer Player
	VSPlayer Player
	Rd       Round
	Status   string
}

func (ann *Announcement) repr() string {
	repr := []string{}
	if ann.Rd != (Round{}) {
		repr = append(repr, "round = "+ann.Rd.ID)
	}
	if ann.VSPlayer != (Player{}) {
		repr = append(repr, "vs = "+ann.VSPlayer.repr())
	}
	repr = append(repr, "status = "+ann.Status)
	return strings.Join(repr, ", ")
}

func (ann *Announcement) toPlayerStatus() *ttt.PlayerStatus {
	ps := ttt.PlayerStatus{}
	ps.RoundID = ann.Rd.ID
	ps.PlayerID = ann.ToPlayer.ID
	ps.PlayerScore = ann.ToPlayer.Score
	if &ann.VSPlayer != nil {
		ps.VSID = ann.VSPlayer.ID
		ps.VSScore = ann.VSPlayer.Score
		ps.VSName = ann.VSPlayer.Name
	}
	ps.Status = ann.Status
	if &ann.Rd != nil {
		ps.GridSnap = ann.Rd.Grid
	}
	return &ps
}

type Group map[string]Round

type TTTServer struct {
	Players       *map[string]*Player
	Groups        *Group
	BenchPlayers  *PlayersQueue
	WithAIPlayers chan *Player

	Announce chan *Announcement // outgoing channel
}

// Create a new round between 2 players.
func (ttts *TTTServer) createNewRound(p1, p2 *Player) Round {
	var grid ttt.Grid

	currentPlayer := p1
	nextPlayer := p2
	if ttt.RandInt(2) == 0 {
		currentPlayer = p2
		nextPlayer = p1
	}
	r := Round{
		ID:            uuid.New(),
		CurrentPlayer: currentPlayer,
		NextPlayer:    nextPlayer,
		Winner:        nil,
		Grid:          &grid,
	}
	currentPlayer.RoundID = r.ID
	nextPlayer.RoundID = r.ID
	(*ttts.Groups)[r.ID] = r
	ttts.Announce <- &Announcement{
		ToPlayer: *r.CurrentPlayer,
		VSPlayer: *r.NextPlayer,
		Rd:       r,
		Status:   ttt.StatusYourTurn,
	}
	ttts.Announce <- &Announcement{
		ToPlayer: *r.NextPlayer,
		VSPlayer: *r.CurrentPlayer,
		Rd:       r,
		Status:   ttt.StatusWaitTurn,
	}
	glog.Infoln("new round between", p1.repr(), "and", p2.repr())
	return r
}

func (ttts *TTTServer) ProcessJoin(p *Player, withAI bool) {
	(*ttts.Players)[p.ID] = p
	glog.Infoln("total players", len((*ttts.Players)))
	ttts.Announce <- &Announcement{
		ToPlayer: *p,
		VSPlayer: Player{},
		Rd:       Round{},
		Status:   ttt.StatusWait,
	}
	if withAI {
		aip := &Player{
			ID:    uuid.New(),
			Name:  BotNames[ttt.RandInt(len(BotNames))],
			Score: ttt.RandInt(100),
		}
		ttts.createNewRound(p, aip)
		glog.Infoln("deploying AI player")
		am.NewAIPlayer(aip.ID)
	} else {
		ttts.BenchPlayers.Remove(p)
		ttts.BenchPlayers.Push(p)
		glog.Infoln("waiting list size", ttts.BenchPlayers.Len())
		if ttts.BenchPlayers.Len() > 1 {
			p1 := ttts.BenchPlayers.Pop()
			p2 := ttts.BenchPlayers.Pop()
			ttts.createNewRound(p1, p2)
		}
	}
}

func (ttts *TTTServer) ProcessQuit(p *Player) {
	delete((*ttts.Players), p.ID)
	if p.RoundID != "" {
		rd := (*ttts.Groups)[p.RoundID]
		// end the round and put the other into waiting queue
		if rd != (Round{}) {
			delete(*ttts.Groups, p.RoundID)
			vs := rd.getOtherPlayer(p)
			ttts.Announce <- &Announcement{
				ToPlayer: *vs,
				VSPlayer: Player{},
				Rd:       Round{},
				Status:   ttt.StatusOtherLeft,
			}
		}
	} else {
		ttts.BenchPlayers.Remove(p)
	}
	glog.Infoln("close connection for player", p.repr())
	if p.WS != nil {
		p.WS.Close()
	}
}

func (ttts *TTTServer) ProcessAnnouncement(a *Announcement) {
	ps := a.toPlayerStatus()
	glog.Infoln("announce to", a.ToPlayer.repr(), ps.Repr())
	if a.ToPlayer.WS != nil {
		a.ToPlayer.WS.WriteJSON(ps)
	} else {
		playerStatuses <- *ps
	}

}

func (ttts *TTTServer) Judge(m *ttt.PlayerAction) {
	rd := (*ttts.Groups)[m.RoundID]
	if rd.ID == "" || rd.CurrentPlayer.ID != m.PlayerID {
		glog.Infoln("Invalid move for player", rd.CurrentPlayer.repr())
		return
	}
	currentUserStatus := ""
	nextUserStatus := ""
	// Switch turn no matter what
	rd.switchTurn()
	if rd.Grid.HasSameMarksInRows(m.Pos, m.PlayerID) {
		rd.Winner = rd.NextPlayer
		rd.CurrentPlayer.Score -= 1
		rd.NextPlayer.Score += 1
		ttts.EndRound(m.RoundID)
		currentUserStatus = ttt.StatusLoss
		nextUserStatus = ttt.StatusWin
	} else if rd.Grid.IsFull() {
		ttts.EndRound(m.RoundID)
		currentUserStatus = ttt.StatusTie
		nextUserStatus = ttt.StatusTie
	} else {
		(*ttts.Groups)[m.RoundID] = rd
		currentUserStatus = ttt.StatusYourTurn
		nextUserStatus = ttt.StatusWaitTurn
	}
	ttts.Announce <- &Announcement{
		ToPlayer: *rd.CurrentPlayer,
		VSPlayer: *rd.NextPlayer,
		Rd:       rd,
		Status:   currentUserStatus,
	}
	ttts.Announce <- &Announcement{
		ToPlayer: *rd.NextPlayer,
		VSPlayer: *rd.CurrentPlayer,
		Rd:       rd,
		Status:   nextUserStatus,
	}
}

func (ttts *TTTServer) EndRound(r string) {
	delete(*ttts.Groups, r)
}

func (ttts *TTTServer) Daemon() {
	for {
		select {
		case a := <-ttts.Announce:
			ttts.ProcessAnnouncement(a)
		case p := <-ttts.WithAIPlayers:
			ttts.ProcessJoin(p, false)
		case a := <-playerActions:
			ttts.Judge(&a)
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
}

func InitTTTServer() *TTTServer {
	ttts := TTTServer{}
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
	go ttts.Daemon()
	return &ttts
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	p := &Player{ws, "", uuid.New(), "", 0}
	ttts.Announce <- &Announcement{
		ToPlayer: *p,
		VSPlayer: Player{},
		Rd:       Round{},
		Status:   ttt.StatusConnected,
	}
	p.parseAction()
}
