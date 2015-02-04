package main

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe"
)

const (
	nameLengthLimit = 8
)

// Fill a range with a give rune.
// x, y: starting position
// w, h: range
func fill(x, y, w, h int, r rune) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, r, ttt.ColDef, ttt.ColDef)
		}
	}
}

// Print lines
// centerX: center x of a terminal
// y: starting position
// msg: messages to be printed
// fg: foreground
// alignCenter: whether or not msg will be center aligned
func printLines(centerX int, y int, msg string, fg termbox.Attribute,
	alignCenter bool) {
	msgs := strings.Split(msg, "\n")
	x := centerX - ttt.Width/2
	for _, m := range msgs {
		if alignCenter {
			x = centerX - len(m)/2
		}
		xstart := x
		for _, c := range m {
			termbox.SetCell(xstart, y, c, fg, ttt.ColDef)
			xstart++
		}
		y++
	}
}

// Get terminal's center
func getTBCenter() ttt.Position {
	w, h := termbox.Size()
	return ttt.Position{w / 2, h / 2}
}

// Convert grid positions to termbox coordinates
func toTBPosition(p ttt.Position) (ttt.Position, error) {
	tbCenter := getTBCenter()
	if !ttt.IsValidPosition(p) {
		return p, errors.New("Invalid position")
	}
	center := ttt.GetCenter()
	x := tbCenter.X - (center.X-p.X)*ttt.Width/ttt.Size
	y := tbCenter.Y - (center.Y-p.Y)*ttt.Height/ttt.Size
	return ttt.Position{x, y}, nil
}

func setCell(p ttt.Position, r rune) {
	tbPos, err := toTBPosition(p)
	if err == nil {
		termbox.SetCell(tbPos.X, tbPos.Y, r, ttt.ColDef, ttt.ColDef)
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
		return ttt.MyRune
	} else if s != "" && s == tttc.VSID {
		return ttt.OtherRune
	} else {
		return ttt.SpecialRune
	}
}

// Check if a cell is available
func (tttc *TTTClient) cellIsPinnable(p ttt.Position) bool {
	return ttt.IsValidPosition(p) && tttc.RoundID != "" &&
		tttc.Grid.Get(p) == ""
}

func (tttc *TTTClient) MoveCursor(direction string) error {
	x := tttc.CursorPos.X
	y := tttc.CursorPos.Y
	switch direction {
	case ttt.Up:
		y--
	case ttt.Down:
		y++
	case ttt.Left:
		x--
	case ttt.Right:
		x++
	}

	if !ttt.IsValidPosition(ttt.Position{x, y}) {
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
	return tttc.ID != "" && tttc.Status == ttt.StatusYourTurn
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
	termbox.Clear(ttt.ColDef, ttt.ColDef)
	tbCenter := getTBCenter()

	tbLeftXPos := tbCenter.X - ttt.Width/2
	tbUpYPos := tbCenter.Y - ttt.Height/2

	// draw the grid
	for yoffset := 0; yoffset <= ttt.Size; yoffset++ {
		for xoffset := 0; xoffset <= ttt.Size; xoffset++ {
			xstart := tbLeftXPos + xoffset*ttt.XSpan
			ystart := tbUpYPos + yoffset*ttt.YSpan
			// all intersections
			termbox.SetCell(xstart, ystart, '+', ttt.ColDef,
				ttt.ColDef)
			if xoffset < ttt.Size {
				fill(xstart+1, ystart, ttt.XSpan-1, 1, '-')
			}
			if yoffset < ttt.Size {
				fill(xstart, ystart+1, 1, ttt.YSpan-1, '|')
			}

		}
	}

	// player score, status, and user manual
	printLines(tbCenter.X, tbUpYPos-2, ttt.Title, ttt.ColDef, true)
	printLines(tbCenter.X, tbUpYPos+ttt.Height+2, tttc.userScores(),
		ttt.ColDef, false)
	printLines(tbCenter.X, tbUpYPos+ttt.Height+4, tttc.Status,
		termbox.ColorBlue, false)
	printLines(tbCenter.X, tbUpYPos+ttt.Height+6, ttt.HelpMsg, ttt.ColDef,
		false)

	tttc.SetCursor(tttc.CursorPos)

	// draw all Xs and Os
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
		RoundID:    tttc.RoundID,
		PlayerID:   tttc.ID,
		PlayerName: tttc.Name,
		Pos:        ttt.Position{},
		Cmd:        cmd,
	}
	return tttc.Conn.WriteJSON(m)
}

func (tttc *TTTClient) SendPin(p ttt.Position) error {
	m := ttt.PlayerAction{
		RoundID:    tttc.RoundID,
		PlayerID:   tttc.ID,
		PlayerName: tttc.Name,
		Pos:        p,
		Cmd:        ttt.CmdMove,
	}
	return tttc.Conn.WriteJSON(m)
}

func (tttc *TTTClient) Update(s ttt.PlayerStatus) error {
	if s.RoundID != "" && tttc.RoundID != s.RoundID &&
		!ttt.IsOverStatus(tttc.Status) {
		glog.Warningln("Round IDs do not match")
		return errors.New("Round IDs do not match")
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

func (tttc *TTTClient) Join(withAI bool) error {
	if !ttt.IsOverStatus(tttc.Status) {
		glog.Warningln("cannot rematch before this round is over")
		return errors.New("This round is not over yet.")
	}
	if withAI {
		return tttc.SendSimpleCMD(ttt.CmdJoinAI)
	}
	return tttc.SendSimpleCMD(ttt.CmdJoin)
}

func (tttc *TTTClient) Quit() error {
	return tttc.SendSimpleCMD(ttt.CmdQuit)
}

// Listen for messages from the server
func (tttc *TTTClient) Listener() error {
	for {
		status := ttt.PlayerStatus{}
		err := tttc.Conn.ReadJSON(&status)
		if err != nil {
			glog.Warningln(tttc.ID, err)
			status = ttt.PlayerStatus{
				RoundID:     tttc.RoundID,
				PlayerName:  tttc.Name,
				PlayerID:    tttc.ID,
				PlayerScore: tttc.Score,
				VSID:        tttc.VSID,
				VSName:      tttc.VSName,
				VSScore:     tttc.VSScore,
				Status:      ttt.StatusLossConnection,
				GridSnap:    &tttc.Grid,
			}
			tttc.Update(status)
			return nil
		}
		tttc.Update(status)

	}
	return nil
}

func TTTCInit(name string) *TTTClient {
	if err := termbox.Init(); err != nil {
		glog.Fatalln(err)
	}
	termbox.SetInputMode(termbox.InputEsc)
	center := ttt.GetCenter()
	tttc := TTTClient{}
	if len(name) > nameLengthLimit {
		tttc.Name = name[:nameLengthLimit]
	} else {
		tttc.Name = name
	}
	tttc.CursorPos = center
	return &tttc
}
