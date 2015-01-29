package server

import (
	"container/list"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe/common"
)

type Player struct {
	WS      *websocket.Conn
	RoundID string
	ID      string
	Name    string
	Score   int
}

// Parse the action sent by a client
func (p *Player) parseAction() {
	for {
		m := ttt.PlayerAction{}
		_ = p.WS.ReadJSON(&m)
		switch m.Cmd {
		case ttt.CMD_QUIT:
			ttts.ProcessQuit(p)
			return
		case ttt.CMD_JOIN:
			p.Name = m.PlayerName
			ttts.ProcessJoin(p, true)
		case ttt.CMD_MOVE:
			ttts.Judge(&m)
		default:
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
	RoundID       string
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
	if r.CurrentPlayer != p {
		return r.CurrentPlayer
	} else if r.NextPlayer != p {
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

func (ann *Announcement) toPlayerStatus() *ttt.PlayerStatus {
	ps := ttt.PlayerStatus{}
	ps.RoundID = ann.Rd.RoundID
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
	Players      *PlayersQueue
	Groups       *Group
	BenchPlayers *PlayersQueue

	Announce chan *Announcement // outgoing channel
}

// Create a new round between 2 players.
func (ttts *TTTServer) createNewRound(p1, p2 *Player) Round {
	var grid ttt.Grid
	rand.Seed(time.Now().Unix())
	currentPlayer := p1
	nextPlayer := p2
	if rand.Intn(2) == 0 {
		currentPlayer = p2
		nextPlayer = p1
	}
	r := Round{
		RoundID:       uuid.New(),
		CurrentPlayer: currentPlayer,
		NextPlayer:    nextPlayer,
		Winner:        nil,
		Grid:          &grid,
	}
	currentPlayer.RoundID = r.RoundID
	nextPlayer.RoundID = r.RoundID
	(*ttts.Groups)[r.RoundID] = r
	ttts.Announce <- &Announcement{*r.CurrentPlayer, *r.NextPlayer,
		r, ttt.STATUS_YOUR_TURN}
	ttts.Announce <- &Announcement{*r.NextPlayer, *r.CurrentPlayer,
		r, ttt.STATUS_WAIT_TURN}
	return r
}

func (ttts *TTTServer) ProcessJoin(p *Player, isNew bool) {
	ttts.BenchPlayers.Push(p)
	ttts.Announce <- &Announcement{
		ToPlayer: *p,
		VSPlayer: Player{},
		Rd:       Round{},
		Status:   ttt.STATUS_WAIT,
	}
	glog.Info("waiting list size ", ttts.BenchPlayers.Len())
	if ttts.BenchPlayers.Len() >= 2 {
		p1 := ttts.BenchPlayers.Pop()
		p2 := ttts.BenchPlayers.Pop()
		ttts.createNewRound(p1, p2)
		glog.Info("New round between ", p1.Name, "(", p1.ID,
			") and ", p2.Name, "(", p2.ID, ")")
	}
	if isNew {
		ttts.Players.Push(p)
	}
	glog.Info("total players ", ttts.Players.Len())
}

func (ttts *TTTServer) ProcessQuit(p *Player) {
	ttts.Players.Remove(p)
	if p.RoundID != "" {
		// end the round and put the other into waiting queue
		rd := (*ttts.Groups)[p.RoundID]
		delete(*ttts.Groups, p.RoundID)
		vs := rd.getOtherPlayer(p)
		vs.RoundID = ""
		ttts.ProcessJoin(vs, false)
	} else {
		ttts.BenchPlayers.Remove(p)
	}
	glog.Info("close connection for player ", p.Name, "(", p.ID, ")")
	p.WS.Close()
}

func (ttts *TTTServer) ProcessAnnouncement(a *Announcement) {
	ps := a.toPlayerStatus()
	glog.Info("announce to ", a.ToPlayer.ID, ps)
	a.ToPlayer.WS.WriteJSON(ps)
}

func (ttts *TTTServer) Judge(m *ttt.PlayerAction) {
	rd := (*ttts.Groups)[m.RoundID]
	if rd.RoundID == "" || rd.CurrentPlayer.ID != m.PlayerID {
		glog.Info("Invalid move for player ", m.PlayerID, rd.RoundID, rd.CurrentPlayer)
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
		currentUserStatus = ttt.STATUS_LOSS
		nextUserStatus = ttt.STATUS_WIN
	} else if rd.Grid.IsFull() {
		ttts.EndRound(m.RoundID)
		currentUserStatus = ttt.STATUS_TIE
		nextUserStatus = ttt.STATUS_TIE
	} else {
		(*ttts.Groups)[m.RoundID] = rd
		currentUserStatus = ttt.STATUS_YOUR_TURN
		nextUserStatus = ttt.STATUS_WAIT_TURN
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

// End this round
func (ttts *TTTServer) EndRound(r string) {
	delete(*ttts.Groups, r)
}

func (ttts *TTTServer) Daemon() {
	for {
		select {
		case a := <-ttts.Announce:
			ttts.ProcessAnnouncement(a)
		default:
		}

	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  READ_BUFFER_SIZE,
	WriteBufferSize: WRITE_BUFFER_SIZE,
}

func Init() *TTTServer {
	ttts := TTTServer{}
	group := make(Group)
	ttts.Players = &PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}

	ttts.BenchPlayers = &PlayersQueue{
		players: list.New(),
		lock:    sync.Mutex{},
	}
	ttts.Announce = make(chan *Announcement, 10)
	ttts.Groups = &group
	go ttts.Daemon()
	return &ttts
}

var ttts = Init()

func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	p := &Player{ws, "", uuid.New(), "", 0}
	p.parseAction()
}
