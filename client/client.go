package client

import (
	"errors"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

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
	PairID    string
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

func PrintLines(x, y int, msg string) {
	xstart := x
	ystart := y
	for _, c := range msg {
		if c == '\n' {
			xstart = x
			ystart++
			continue
		}
		termbox.SetCell(xstart, ystart, c, ttt.COLDEF, ttt.COLDEF)
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
	if s == tttc.Name {
		return ttt.MYRUNE
	} else if s == tttc.VSName {
		return ttt.OTHERRUNE
	} else {
		return ttt.SPECIALRUNE
	}
}

// Check if a cell is available
func (tttc *TTTClient) CellIsAvailable(p ttt.Position) bool {
	return tttc.Grid.Get(p) == ""
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

func (tttc *TTTClient) PinCursor(r rune) bool {
	if tttc.CellIsAvailable(tttc.CursorPos) {
		tttc.Grid.Set(tttc.CursorPos, tttc.Name)
		return true
	} else {
		return false
	}
}

func (tttc *TTTClient) DrawCells() {
	for x := 0; x < ttt.SIZE; x++ {
		for y := 0; y < ttt.SIZE; y++ {
			p := ttt.Position{x, y}
			r := tttc.NameToRune(tttc.Grid.Get(p))
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

	PrintLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+1, tttc.Status)
	PrintLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+3, ttt.HELPMSG)

	tttc.SetCursor(tttc.CursorPos)
	// draw all on cells
	tttc.DrawCells()
	termbox.Flush()
}

func Init() *TTTClient {
	if err := termbox.Init(); err != nil {
		panic(err)
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
	tttc.Status = WAIT_STATUS
	return nil
}

func (tttc *TTTClient) SendPayload() error {

}
