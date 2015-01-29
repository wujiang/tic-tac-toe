package client

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe/common"
)

const (
	NAME_LENGTH_LIMIT = 8
)

func fill(x, y, w, h int, r rune) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, r, ttt.COLDEF, ttt.COLDEF)
		}
	}
}

func printLines(x, y int, msg string, fg termbox.Attribute) {
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

func getTBCenter() ttt.Position {
	w, h := termbox.Size()
	return ttt.Position{w / 2, h / 2}
}

func getCenter() ttt.Position {
	x := (ttt.SIZE - 1) / 2
	y := (ttt.SIZE - 1) / 2
	return ttt.Position{x, y}
}

func isValidatePosition(p ttt.Position) bool {
	return p.X >= 0 && p.X < ttt.SIZE && p.Y >= 0 && p.Y < ttt.SIZE
}

// Convert grid positions to termbox coordinates
func toTBPosition(p ttt.Position) (ttt.Position, error) {
	tbCenter := getTBCenter()
	if !isValidatePosition(p) {
		return p, errors.New("Invalid position")
	}
	center := getCenter()
	x := tbCenter.X - (center.X-p.X)*ttt.WIDTH/ttt.SIZE
	y := tbCenter.Y - (center.Y-p.Y)*ttt.HEIGHT/ttt.SIZE
	return ttt.Position{x, y}, nil
}

func setCell(p ttt.Position, r rune) {
	tbPos, err := toTBPosition(p)
	if err == nil {
		termbox.SetCell(tbPos.X, tbPos.Y, r, ttt.COLDEF, ttt.COLDEF)
	}

}

type TTTClient struct {
	Name      string
	ID        string
	Score     int
	Conn      *websocket.Conn
	VSID      string
	VSName    string
	VSScore   int
	RoundID   string
	Status    string
	CursorPos ttt.Position
	Grid      ttt.Grid
}

func (tttc *TTTClient) nameToRune(s string) rune {
	if s != "" && s == tttc.ID {
		return ttt.MYRUNE
	} else if s != "" && s == tttc.VSID {
		return ttt.OTHERRUNE
	} else {
		return ttt.SPECIALRUNE
	}
}

// Check if a cell is available
func (tttc *TTTClient) cellIsPinnable(p ttt.Position) bool {
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

	if !isValidatePosition(ttt.Position{x, y}) {
		return errors.New("Invalid position")
	}

	tttc.CursorPos.X = x
	tttc.CursorPos.Y = y
	return nil
}

func (tttc *TTTClient) SetCursor(p ttt.Position) error {
	tbPos, err := toTBPosition(p)
	if err == nil {
		termbox.SetCursor(tbPos.X, tbPos.Y)
		return nil
	}
	return err
}

func (tttc *TTTClient) isYourTurn() bool {
	return tttc.ID != "" && tttc.Status == ttt.STATUS_YOUR_TURN
}

func (tttc *TTTClient) PinCursor(r rune) bool {
	if tttc.cellIsPinnable(tttc.CursorPos) && tttc.isYourTurn() {
		tttc.Grid.Set(tttc.CursorPos, tttc.ID)
		err := tttc.SendPin(tttc.CursorPos)
		return err == nil
	} else {
		return false
	}
}

func (tttc *TTTClient) drawCells() {
	for x, l := range tttc.Grid {
		for y, s := range l {
			p := ttt.Position{x, y}
			r := tttc.nameToRune(s)
			setCell(p, r)
		}
	}
}

func (tttc *TTTClient) userScores() string {
	var buffer bytes.Buffer
	buffer.WriteString(tttc.Name)
	buffer.WriteString(": ")
	buffer.WriteString(strconv.Itoa(tttc.Score))
	if tttc.VSName != "" {
		buffer.WriteString(" VS ")
		buffer.WriteString(tttc.VSName)
		buffer.WriteString(": ")
		buffer.WriteString(strconv.Itoa(tttc.VSScore))
	}
	return buffer.String()
}

func (tttc *TTTClient) RedrawAll() {
	termbox.Clear(ttt.COLDEF, ttt.COLDEF)
	tbCenter := getTBCenter()

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
				fill(xstart+1, ystart, ttt.XSPAN-1, 1, '-')
			}
			if yoffset < ttt.SIZE {
				fill(xstart, ystart+1, 1, ttt.YSPAN-1, '|')
			}

		}
	}

	printLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+2, tttc.userScores(),
		ttt.COLDEF)
	printLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+4, tttc.Status,
		termbox.ColorBlue)
	printLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+6, ttt.HELPMSG, ttt.COLDEF)

	tttc.SetCursor(tttc.CursorPos)
	// draw all on cells
	tttc.drawCells()
	termbox.Flush()
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
		tttc.ID,
		tttc.Name,
		ttt.Position{},
		cmd,
	}
	return tttc.Conn.WriteJSON(m)
}

func (tttc *TTTClient) SendPin(p ttt.Position) error {
	m := ttt.PlayerAction{
		tttc.RoundID,
		tttc.ID,
		tttc.Name,
		p,
		ttt.CMD_MOVE,
	}
	return tttc.Conn.WriteJSON(m)
}

func (tttc *TTTClient) Update(s ttt.PlayerStatus) error {
	if s.RoundID != "" && tttc.RoundID != s.RoundID &&
		!ttt.IsOverStatus(tttc.Status) {
		glog.Warning("Round IDs do not match")
	} else {
		tttc.RoundID = s.RoundID
	}
	tttc.ID = s.PlayerID
	tttc.Score = s.PlayerScore
	tttc.VSID = s.VSID
	tttc.VSName = s.VSName
	tttc.VSScore = s.VSScore
	tttc.Status = s.Status
	if s.GridSnap != nil {
		tttc.Grid = *s.GridSnap
	} else {
		var grid ttt.Grid
		tttc.Grid = grid
	}
	tttc.RedrawAll()
	return nil
}

func (tttc *TTTClient) Join() error {
	if !ttt.IsOverStatus(tttc.Status) {
		return errors.New("This round is not over yet.")
	}
	err := tttc.SendSimpleCMD(ttt.CMD_JOIN)
	return err
}

func (tttc *TTTClient) Listener() error {
	for {
		status := ttt.PlayerStatus{}
		err := tttc.Conn.ReadJSON(&status)
		if err != nil {
			glog.Warning(tttc.ID, ", ", err)
			status = ttt.PlayerStatus{
				tttc.RoundID,
				tttc.Name,
				tttc.ID,
				tttc.Score,
				tttc.VSID,
				tttc.VSName,
				tttc.VSScore,
				ttt.STATUS_LOSS_CONNECTION,
				&tttc.Grid,
			}
			tttc.Update(status)
			return nil
		}
		tttc.Update(status)

	}
	return nil
}

func Init(name string) *TTTClient {
	if err := termbox.Init(); err != nil {
		glog.Fatal(err)
	}

	termbox.SetInputMode(termbox.InputEsc)
	center := getCenter()
	tttc := TTTClient{}
	if len(name) > NAME_LENGTH_LIMIT {
		tttc.Name = name[:NAME_LENGTH_LIMIT]
	} else {
		tttc.Name = name
	}
	tttc.CursorPos = center
	return &tttc
}
