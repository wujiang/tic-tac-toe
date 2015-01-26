package server

import (
	"container/list"
	"net/http"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe/common"
)

type Player struct {
	WS      *websocket.Conn
	RoundID string
	Name    string
}

// Parse the action sent by a client
func (p *Player) parseAction() {
	for {
		m := ttt.PlayerAction{}
		_ = p.WS.ReadJSON(&m)

		switch m.Cmd {
		case ttt.CMD_QUIT:
			ttts.ProcessQuit(p)
			break
		case ttt.CMD_JOIN:
			ttts.ProcessJoin(p)
		case ttt.CMD_MOVE:
			ttts.Judge(&m)
		}
	}

	glog.Info("close connection for player ", p.Name)
	p.WS.Close()

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

type Group map[string]Round

type Announcement struct {
	ToPlayer Player
	VSPlayer Player
	Rd       Round
	Status   string
}

func (ann *Announcement) toPlayerStatus() *ttt.PlayerStatus {
	ps := ttt.PlayerStatus{}
	ps.RoundID = ann.Rd.RoundID
	ps.PlayerName = ann.ToPlayer.Name
	if &ann.VSPlayer != nil {
		ps.VSName = ann.VSPlayer.Name
	}
	ps.Status = ann.Status
	if &ann.Rd != nil {
		ps.GridSnap = ann.Rd.Grid
	}
	return &ps
}

type TTTServer struct {
	Groups       *Group
	BenchPlayers *list.List

	Join     chan *Player       // incoming channel
	Quit     chan *Player       // incoming channel
	Announce chan *Announcement // outgoing channel
}

func (ttts *TTTServer) ProcessJoin(p *Player) {
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	// put into the waiting queue
	ttts.BenchPlayers.PushBack(p)
	glog.Info("push ", p.Name, " into waiting list")
	glog.Info("waiting list size ", ttts.BenchPlayers.Len())
	if ttts.BenchPlayers.Len() <= 1 {
		ttts.Announce <- &Announcement{*p, Player{}, Round{},
			ttt.STATUS_WAIT}
	} else {
		r := Round{}
		r.RoundID = uuid.New()
		p1 := ttts.BenchPlayers.Remove(ttts.BenchPlayers.Front())
		p2 := ttts.BenchPlayers.Remove(ttts.BenchPlayers.Front())
		r.CurrentPlayer = p1.(*Player)
		r.NextPlayer = p2.(*Player)
		var grid ttt.Grid
		r.Grid = &grid
		(*ttts.Groups)[r.RoundID] = r
		ttts.Announce <- &Announcement{*r.CurrentPlayer, *r.NextPlayer, r,
			ttt.STATUS_YOUR_TURN}
		ttts.Announce <- &Announcement{*r.NextPlayer, *r.CurrentPlayer, r,
			ttt.STATUS_WAIT_TURN}

	}
}

func (ttts *TTTServer) ProcessQuit(p *Player) {
	if p.RoundID != "" {
		// end the round and put the other into waiting queue
		rd := (*ttts.Groups)[p.RoundID]
		delete(*ttts.Groups, p.RoundID)
		vs := rd.getOtherPlayer(p)
		vs.RoundID = ""
		ttts.ProcessJoin(vs)
	} else {
		lock := sync.Mutex{}
		lock.Lock()
		for e := ttts.BenchPlayers.Front(); e != nil; e = e.Next() {
			if e.Value.(*Player) == p {
				ttts.BenchPlayers.Remove(e)
				break
			}
		}
		lock.Unlock()

	}
}

func (ttts *TTTServer) ProcessAnnouncement(a *Announcement) {
	ps := a.toPlayerStatus()
	glog.Info("announce to ", a.ToPlayer.Name, ": ", ps)
	a.ToPlayer.WS.WriteJSON(ps)
}

func (ttts *TTTServer) Judge(m *ttt.PlayerAction) {
	rd := (*ttts.Groups)[m.RoundID]
	if rd.RoundID == "" || rd.CurrentPlayer.Name != m.PlayerName {
		glog.Info("Invalid move for player ", m.PlayerName)
		return
	}
	currentUserStatus := ""
	nextUserStatus := ""

	if rd.Grid.HasSameMarksInRows(m.Pos, m.PlayerName) {
		rd.Winner = rd.CurrentPlayer
		ttts.EndRound(m.RoundID)
		currentUserStatus = ttt.STATUS_WIN
		nextUserStatus = ttt.STATUS_LOSS
	} else if rd.Grid.IsFull() {
		ttts.EndRound(m.RoundID)
		currentUserStatus = ttt.STATUS_TIE
		nextUserStatus = ttt.STATUS_TIE

	} else {
		rd.switchTurn()
		(*ttts.Groups)[m.RoundID] = rd
		currentUserStatus = ttt.STATUS_YOUR_TURN
		nextUserStatus = ttt.STATUS_WAIT_TURN
	}
	ttts.Announce <- &Announcement{*rd.CurrentPlayer, *rd.NextPlayer,
		rd, currentUserStatus}
	ttts.Announce <- &Announcement{*rd.NextPlayer, *rd.CurrentPlayer,
		rd, nextUserStatus}

}

// End this round
// - Remove round from server's rounds list
// - Do not add either user into the waiting queue (should be done in rejoin game)
// - Do not close connections
func (ttts *TTTServer) EndRound(r string) {
	delete(*ttts.Groups, r)
}

func (ttts *TTTServer) Daemon() {
	for {
		select {
		case j := <-ttts.Join:
			ttts.ProcessJoin(j)
		case a := <-ttts.Announce:
			ttts.ProcessAnnouncement(a)
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
	ttts.BenchPlayers = list.New()
	ttts.Join = make(chan *Player, 10)
	ttts.Quit = make(chan *Player, 10)
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
	p := &Player{ws, "", uuid.New()}
	ttts.Join <- p
	defer func() {
		ttts.Quit <- p
	}()

	p.parseAction()
}
