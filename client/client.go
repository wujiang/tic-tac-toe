package client

import (
	"errors"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe/common"
)

const (
	WAIT_STATUS string = "Waiting for another player"
)

type TTTClient struct {
	Name      string
	Conn      *websocket.Conn
	VSName    string
	RoundID   string
	Status    string
	CursorPos ttt.Position // cursor position
	Grid      ttt.Grid
}

func Fill(x, y, w, h int, r rune) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, r, ttt.COLDEF, ttt.COLDEF)
		}
	}
}

func PrintLines(x, y int, msg string, fg termbox.Attribute) {
	xstart := x
	ystart := y
	for _, c := range msg {
		if c == '\n' {
			xstart = x
			ystart++
			continue
		}
		termbox.SetCell(xstart, ystart, c, fg, ttt.COLDEF)
		xstart++
	}

}

func GetTBCenter() ttt.Position {
	w, h := termbox.Size()
	return ttt.Position{w / 2, h / 2}
}

func GetCenter() ttt.Position {
	x := (ttt.SIZE - 1) / 2
	y := (ttt.SIZE - 1) / 2
	return ttt.Position{x, y}
}

func IsValidatePosition(p ttt.Position) bool {
	return p.X >= 0 && p.X < ttt.SIZE && p.Y >= 0 && p.Y < ttt.SIZE
}

// Convert grid positions to termbox coordinates
func ToTBPosition(p ttt.Position) (ttt.Position, error) {
	tbCenter := GetTBCenter()
	if !IsValidatePosition(p) {
		return p, errors.New("Invalid position")
	}
	center := GetCenter()
	x := tbCenter.X - (center.X-p.X)*ttt.WIDTH/ttt.SIZE
	y := tbCenter.Y - (center.Y-p.Y)*ttt.HEIGHT/ttt.SIZE
	return ttt.Position{x, y}, nil
}

func SetCell(p ttt.Position, r rune) {
	tbPos, err := ToTBPosition(p)
	if err == nil {
		termbox.SetCell(tbPos.X, tbPos.Y, r, ttt.COLDEF, ttt.COLDEF)
	}

}

func (tttc *TTTClient) NameToRune(s string) rune {
	if s != "" && s == tttc.Name {
		return ttt.MYRUNE
	} else if s != "" && s == tttc.VSName {
		return ttt.OTHERRUNE
	} else {
		return ttt.SPECIALRUNE
	}
}

// Check if a cell is available
func (tttc *TTTClient) CellIsPinnable(p ttt.Position) bool {
	return tttc.RoundID != "" && tttc.Grid.Get(p) == ""
}

func (tttc *TTTClient) MoveCursor(direction string) error {
	x := tttc.CursorPos.X
	y := tttc.CursorPos.Y
	switch direction {
	case ttt.UP:
		y--
	case ttt.DOWN:
		y++
	case ttt.LEFT:
		x--
	case ttt.RIGHT:
		x++
	}

	if !IsValidatePosition(ttt.Position{x, y}) {
		return errors.New("Invalid position")
	}

	tttc.CursorPos.X = x
	tttc.CursorPos.Y = y
	return nil
}

func (tttc *TTTClient) SetCursor(p ttt.Position) error {
	tbPos, err := ToTBPosition(p)
	if err == nil {
		termbox.SetCursor(tbPos.X, tbPos.Y)
		return nil
	}
	return err
}

func (tttc *TTTClient) isYourTurn() bool {
	return tttc.Name != "" && tttc.Status == ttt.STATUS_YOUR_TURN
}

func (tttc *TTTClient) PinCursor(r rune) bool {
	if tttc.CellIsPinnable(tttc.CursorPos) && tttc.isYourTurn() {
		tttc.Grid.Set(tttc.CursorPos, tttc.Name)
		err := tttc.SendPin(tttc.CursorPos)
		glog.Info("pincursor err", err)
		return err == nil
	} else {
		return false
	}
}

func (tttc *TTTClient) DrawCells() {
	for x, l := range tttc.Grid {
		for y, s := range l {
			p := ttt.Position{x, y}
			r := tttc.NameToRune(s)
			SetCell(p, r)
		}
	}
}

func (tttc *TTTClient) RedrawAll() {
	termbox.Clear(ttt.COLDEF, ttt.COLDEF)
	tbCenter := GetTBCenter()

	tbLeftXPos := tbCenter.X - ttt.WIDTH/2
	tbUpYPos := tbCenter.Y - ttt.HEIGHT/2

	for yoffset := 0; yoffset <= ttt.SIZE; yoffset++ {
		for xoffset := 0; xoffset <= ttt.SIZE; xoffset++ {
			xstart := tbLeftXPos + xoffset*ttt.XSPAN
			ystart := tbUpYPos + yoffset*ttt.YSPAN
			// all intersections
			termbox.SetCell(xstart, ystart, '+', ttt.COLDEF,
				ttt.COLDEF)
			if xoffset < ttt.SIZE {
				Fill(xstart+1, ystart, ttt.XSPAN-1, 1, '-')
			}
			if yoffset < ttt.SIZE {
				Fill(xstart, ystart+1, 1, ttt.YSPAN-1, '|')
			}

		}
	}

	PrintLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+1, tttc.Name,
		termbox.ColorYellow)
	PrintLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+3, tttc.Status,
		termbox.ColorBlue)
	PrintLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+5, ttt.HELPMSG, ttt.COLDEF)

	tttc.SetCursor(tttc.CursorPos)
	// draw all on cells
	tttc.DrawCells()
	termbox.Flush()
}

func Init() *TTTClient {
	if err := termbox.Init(); err != nil {
		glog.Fatal(err)
	}

	termbox.SetInputMode(termbox.InputEsc)
	center := GetCenter()
	tttc := TTTClient{}
	tttc.Name = uuid.New()

	// TODO: set on connect
	tttc.VSName = uuid.New()

	tttc.CursorPos = center
	return &tttc
}

func (tttc *TTTClient) Connect(s string) error {
	dialer := websocket.DefaultDialer
	ws, _, err := dialer.Dial(s, http.Header{})
	if err != nil {
		return err
	}
	tttc.Conn = ws
	return nil
}

func (tttc *TTTClient) SendSimpleCMD(cmd string) error {
	m := ttt.PlayerAction{
		tttc.RoundID,
		tttc.Name,
		ttt.Position{},
		cmd,
	}
	glog.Info("sending to server:", m)
	return tttc.Conn.WriteJSON(m)
}

func (tttc *TTTClient) SendPin(p ttt.Position) error {
	m := ttt.PlayerAction{
		tttc.RoundID,
		tttc.Name,
		p,
		ttt.CMD_MOVE,
	}
	glog.Info("sending to server:", m)
	return tttc.Conn.WriteJSON(m)
}

func (tttc *TTTClient) Update(s ttt.PlayerStatus) error {
	if tttc.RoundID == "" {
		tttc.RoundID = s.RoundID
	} else if tttc.RoundID != s.RoundID {
		return errors.New("Round IDs do not match")
	}
	glog.Info("s", s)
	tttc.Name = s.PlayerName
	tttc.VSName = s.VSName
	tttc.Status = s.Status
	if s.GridSnap != nil {
		tttc.Grid = *s.GridSnap
	}
	tttc.RedrawAll()
	return nil
}

func (tttc *TTTClient) Listener() error {
	for {
		status := ttt.PlayerStatus{}
		err := tttc.Conn.ReadJSON(&status)
		if err != nil {
			glog.Info(err)
		}
		tttc.Update(status)

	}
	return nil
}
