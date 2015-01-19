package client

import (
	"errors"

	"github.com/nsf/termbox-go"
)

type Position struct {
	x int
	y int
}

type TBCell struct {
	TBPos Position // termbox position
	ch    rune
}

type TTTClient struct {
	cursorPos Position // cursor position
	grid      *map[Position]TBCell
}

func (tbc *TBCell) DrawCell() {
	if tbc.ch != SPECIALRUNE {
		termbox.SetCell(tbc.TBPos.x, tbc.TBPos.y, tbc.ch, COLDEF, COLDEF)
	}
}

func Fill(x, y, w, h int, r rune) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, r, COLDEF, COLDEF)
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
		termbox.SetCell(xstart, ystart, c, COLDEF, COLDEF)
		xstart++
	}

}

func GetTBCenter() Position {
	w, h := termbox.Size()
	return Position{w / 2, h / 2}
}

func GetCenter() Position {
	x := (SIZE-1)/2 + 1
	y := (SIZE-1)/2 + 1
	return Position{x, y}
}

func ValidatePosition(p Position) error {
	if p.x < 1 || p.x > SIZE || p.y < 1 || p.y > SIZE {
		return errors.New("Positions are out of boundary")
	}
	return nil
}

// Convert grid positions to termbox coordinates
func ToTBPosition(p Position, tbCenter Position) (Position, error) {
	if err := ValidatePosition(p); err != nil {
		return p, err
	}
	center := GetCenter()
	x := tbCenter.x - (center.x-p.x)*WIDTH/SIZE
	y := tbCenter.y - (center.y-p.y)*HEIGHT/SIZE
	return Position{x, y}, nil
}

func (ttt *TTTClient) MoveCursor(direction string) error {
	x := ttt.cursorPos.x
	y := ttt.cursorPos.y
	switch direction {
	case UP:
		y--
	case DOWN:
		y++
	case LEFT:
		x--
	case RIGHT:
		x++
	}

	if err := ValidatePosition(Position{x, y}); err != nil {
		return err
	}

	ttt.cursorPos.x = x
	ttt.cursorPos.y = y
	return nil
}

func (ttt *TTTClient) SetCursor(p Position) error {
	tbcell := (*ttt.grid)[p]
	if &tbcell != nil {
		termbox.SetCursor(tbcell.TBPos.x, tbcell.TBPos.y)
		return nil
	}
	return errors.New("Invalid position")
}

func (ttt *TTTClient) PinCursor(r rune) {
	tbPos := (*ttt.grid)[ttt.cursorPos].TBPos
	(*ttt.grid)[ttt.cursorPos] = TBCell{tbPos, r}
}

func (ttt *TTTClient) RedrawAll() {
	termbox.Clear(COLDEF, COLDEF)
	tbCenter := GetTBCenter()

	tbLeftXPos := tbCenter.x - WIDTH/2
	tbUpYPos := tbCenter.y - HEIGHT/2

	for yoffset := 0; yoffset <= SIZE; yoffset++ {
		for xoffset := 0; xoffset <= SIZE; xoffset++ {
			xstart := tbLeftXPos + xoffset*XSPAN
			ystart := tbUpYPos + yoffset*YSPAN
			// all intersections
			termbox.SetCell(xstart, ystart, '+', COLDEF, COLDEF)
			if xoffset < SIZE {
				Fill(xstart+1, ystart, XSPAN-1, 1, '-')
			}
			if yoffset < SIZE {
				Fill(xstart, ystart+1, 1, YSPAN-1, '|')
			}

		}
	}

	PrintLines(tbLeftXPos, tbUpYPos+HEIGHT+1, HELPMSG)

	ttt.SetCursor(ttt.cursorPos)
	// draw all on cells
	for _, v := range *ttt.grid {
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
	tttc.cursorPos = center
	grid := make(map[Position]TBCell)
	tttc.grid = &grid
	for x := 1; x <= SIZE; x++ {
		for y := 1; y <= SIZE; y++ {
			p := Position{x, y}
			tbPos, _ := ToTBPosition(p, tbCenter)
			grid[p] = TBCell{tbPos, SPECIALRUNE}
		}
	}
	return &tttc
}
