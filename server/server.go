package server

import (
	"container/list"
	"fmt"
	"net/http"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe/common"
)

type Player struct {
	WS      *websocket.Conn
	RoundID string
	Name    string
}

// Parse the action sent by a client
func (p *Player) ParseAction() {
	for {
		m := ttt.PlayerAction{}
		err := p.WS.ReadJSON(&m)
		fmt.Println(m, err)
		if err != nil {
			break
		}

		switch m.Cmd {
		case ttt.CMD_QUIT:
			fmt.Println("parse quit")
			break
		case ttt.CMD_JOIN:
			fmt.Println("parse join")
			ttts.ProcessJoin(p)
		case ttt.CMD_MOVE:
			fmt.Println("parse move")
			ttts.Judge(&m)
		}
	}

	fmt.Println("close connection")
	p.WS.Close()

}

func (p *Player) Act() {

}

type Round struct {
	RoundID       string
	CurrentPlayer *Player
	NextPlayer    *Player
	Winner        *Player
	Grid          *ttt.Grid
}

// Switch turn in a matching round
func (r *Round) SwitchTurn() {
	temp := r.CurrentPlayer
	r.CurrentPlayer = r.NextPlayer
	r.NextPlayer = temp
}

type Group map[string]Round

type Announcement struct {
	ToPlayer Player
	VSPlayer Player
	Rd       Round
	Status   string
}

func (ann *Announcement) ToPlayerStatus() *ttt.PlayerStatus {
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

	Join        chan *Player           // incoming channel
	Quit        chan *Player           // incoming channel
	PlayerMoves chan *ttt.PlayerAction // incoming channel
	Announce    chan *Announcement     // outgoing channel
}

// func (r *Round) GetPlayer(p string) *Player {
// 	if r.P1.Name == p {
// 		return r.P1
// 	} else if r.P2.Name == p {
// 		return r.P2
// 	} else {
// 		return nil
// 	}
// }

func (ttts *TTTServer) ProcessJoin(p *Player) {
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	// put into the waiting queue
	ttts.BenchPlayers.PushBack(p)
	fmt.Println("push", p.Name, *p, "into waiting list")
	fmt.Println("waiting list", ttts.BenchPlayers.Len())
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
		(*ttts.Groups)[r.RoundID] = r
		ann1 := Announcement{*r.CurrentPlayer, *r.NextPlayer, r,
			ttt.STATUS_YOUR_TURN}
		ann2 := Announcement{*r.NextPlayer, *r.CurrentPlayer, r,
			ttt.STATUS_WAIT}
		ttts.Announce <- &ann1
		ttts.Announce <- &ann2
	}
}

func (ttts *TTTServer) ProcessAnnouncement(a *Announcement) {
	ps := a.ToPlayerStatus()
	a.ToPlayer.WS.WriteJSON(ps)
}

// func (ttts *TTTServer) QuitWaiting(p *Player) {
// 	lock := sync.Mutex{}
// 	lock.Lock()
// 	defer lock.Unlock()
// 	for e := ttts.BenchPlayers.Front(); e != nil; e = e.Next() {
// 		if e.Value == p {
// 			ttts.BenchPlayers.Remove(e)
// 			break
// 		}
// 	}
// }

// // Remove round from Groups and put the other user into BenchPlayers
// func (ttts *TTTServer) QuitRound(m *ttt.PlayerAction) {
// 	rd := (*ttts.Groups)[m.RoundID]
// 	if &rd == nil {
// 		return
// 	}
// 	toNotify := rd.P1
// 	if rd.P1.Name == m.PlayerName {
// 		toNotify = rd.P2
// 		rd.P1 = nil
// 	} else {
// 		rd.P2 = nil
// 	}
// 	// ttts.Rounds <- &rd

// 	// TODO: remove round from groups

// 	ttts.ProcessJoin(toNotify)
// }

// TODO: ttts.Rounds may have a nil player, which means the other player is gone

func (ttts *TTTServer) Judge(m *ttt.PlayerAction) {
	rd := (*ttts.Groups)[m.RoundID]
	if rd.CurrentPlayer.Name != m.PlayerName {
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
		rd.SwitchTurn()
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
		case m := <-ttts.PlayerMoves:
			fmt.Println("process move")
			ttts.Judge(m)
		case j := <-ttts.Join:
			fmt.Println("process join")
			ttts.ProcessJoin(j)
		case a := <-ttts.Announce:
			fmt.Println("process announcement")
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
	ttts.Join = make(chan *Player)
	ttts.Quit = make(chan *Player)
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
	fmt.Println("Incoming connection", *p)
	ttts.Join <- p
	defer func() {
		ttts.Quit <- p
	}()

	go p.Act()

	p.ParseAction()
}
