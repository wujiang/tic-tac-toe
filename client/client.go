package client

import (
	"errors"

	"code.google.com/p/go-uuid/uuid"

	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe/common"
)

type TBCell struct {
	TBPos ttt.Position // termbox position
	Ch    rune
}

type TTTClient struct {
	Name      string
	CursorPos ttt.Position // cursor position
	Grid      *map[ttt.Position]TBCell
}

func (tbc *TBCell) DrawCell() {
	if tbc.Ch != ttt.SPECIALRUNE {
		termbox.SetCell(tbc.TBPos.X, tbc.TBPos.Y, tbc.Ch, ttt.COLDEF,
			ttt.COLDEF)
	}
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
	x := (ttt.SIZE-1)/2 + 1
	y := (ttt.SIZE-1)/2 + 1
	return ttt.Position{x, y}
}

func IsValidatePosition(p ttt.Position) bool {
	return p.X >= 1 && p.X <= ttt.SIZE && p.Y >= 1 && p.Y <= ttt.SIZE
}

// Convert grid positions to termbox coordinates
func ToTBPosition(p ttt.Position, tbCenter ttt.Position) (ttt.Position, error) {
	if !IsValidatePosition(p) {
		return p, errors.New("Invalid position")
	}
	center := GetCenter()
	x := tbCenter.X - (center.X-p.X)*ttt.WIDTH/ttt.SIZE
	y := tbCenter.Y - (center.Y-p.Y)*ttt.HEIGHT/ttt.SIZE
	return ttt.Position{x, y}, nil
}

// Check if a cell is available
func (tttc *TTTClient) CellIsAvailable(p ttt.Position) bool {
	tbcell := (*tttc.Grid)[p]
	if &tbcell != nil {
		return tbcell.Ch == ttt.SPECIALRUNE
	}
	return false
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
	tbcell := (*tttc.Grid)[p]
	if &tbcell != nil {
		termbox.SetCursor(tbcell.TBPos.X, tbcell.TBPos.Y)
		return nil
	}
	return errors.New("Invalid position")
}

func (tttc *TTTClient) PinCursor(r rune) bool {
	if tttc.CellIsAvailable(tttc.CursorPos) {
		tbPos := (*tttc.Grid)[tttc.CursorPos].TBPos
		(*tttc.Grid)[tttc.CursorPos] = TBCell{tbPos, r}
		return true
	} else {
		return false
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

	PrintLines(tbLeftXPos, tbUpYPos+ttt.HEIGHT+3, ttt.HELPMSG)

	tttc.SetCursor(tttc.CursorPos)
	// draw all on cells
	for _, v := range *tttc.Grid {
		v.DrawCell()
	}
	termbox.Flush()
}

func Init() *TTTClient {
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	termbox.SetInputMode(termbox.InputEsc)
	tbCenter := GetTBCenter()
	center := GetCenter()
	tttc := TTTClient{}
	tttc.Name = uuid.New()
	tttc.CursorPos = center
	grid := make(map[ttt.Position]TBCell)
	tttc.Grid = &grid
	for x := 1; x <= ttt.SIZE; x++ {
		for y := 1; y <= ttt.SIZE; y++ {
			p := ttt.Position{x, y}
			tbPos, _ := ToTBPosition(p, tbCenter)
			grid[p] = TBCell{tbPos, ttt.SPECIALRUNE}
		}
	}
	return &tttc
}
