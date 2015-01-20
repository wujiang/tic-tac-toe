package server

import (
	"container/list"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe/common"
)

type Conn struct {
	ws   *websocket.Conn
	send chan []byte
}

type Player struct {
	Name string
	Cn   Conn
}

type Pair struct {
	P1       *Player
	P2       *Player
	NextMove *Player
}

type Grid map[ttt.Position]string

type GridStatus struct {
	PairID string
	Winner *Player
	Grd    *Grid
}

type PlayerMove struct {
	PairID string
	Plyer  *Player
	Cll    ttt.Position
}

type TTTServer struct {
	Groups       *map[string]Pair
	BenchPlayers *list.List

	PlayerMove chan []PlayerMove
	GrdSts     chan []GridStatus
}

func (ttts *TTTServer) Judge(m PlayerMove) {

}

func (ttts *TTTServer) run() {
	for {
		select {
		case m := <-ttts.PlayerMove:
			ttts.GrdSts <- ttts.Judge(m)
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  READ_BUFFER_SIZE,
	WriteBufferSize: WRITE_BUFFER_SIZE,
}

var palyersQueue = list.New()

func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
}
