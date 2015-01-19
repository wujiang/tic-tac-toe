package client

import "github.com/nsf/termbox-go"

const WIDTH int = 30
const HEIGHT int = 12
const SIZE int = 3
const XSPAN int = WIDTH / SIZE
const YSPAN int = HEIGHT / SIZE
const COLDEF = termbox.ColorDefault
const UP string = "up"
const DOWN string = "down"
const LEFT string = "left"
const RIGHT string = "right"
const SPECIALRUNE rune = '_'
const MYRUNE rune = 'O'
const OTHERRUNE rune = 'X'
const HELPMSG = `
Tic-tac-toe manual:
- LEFT: h, ctrl-b, arrow-left
- DOWN: j, ctrl-n, arrow-down
- UP: k, ctrl-p, arrow-up
- RIGHT: l, ctrl-f, arrow-right
- EXIT: q, esc
- ENTER: o, enter
`
