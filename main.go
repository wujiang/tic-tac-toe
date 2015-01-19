package main

import (
	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe/client"
)

func main() {
	tttc := client.Init()
	defer termbox.Close()

	tttc.RedrawAll()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			// arrows, enter, and esc
			switch ev.Key {
			case termbox.KeyEnter:
				tttc.PinCursor(client.MYRUNE)
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft:
				tttc.MoveCursor(client.LEFT)
			case termbox.KeyArrowDown:
				tttc.MoveCursor(client.DOWN)
			case termbox.KeyArrowUp:
				tttc.MoveCursor(client.UP)
			case termbox.KeyArrowRight:
				tttc.MoveCursor(client.RIGHT)

			}

			// vim binding
			switch ev.Ch {
			case 'o':
				tttc.PinCursor(client.MYRUNE)
			case 'q':
				break mainloop
			case 'h':
				tttc.MoveCursor(client.LEFT)
			case 'j':
				tttc.MoveCursor(client.DOWN)
			case 'k':
				tttc.MoveCursor(client.UP)
			case 'l':
				tttc.MoveCursor(client.RIGHT)
			}

		case termbox.EventError:
			panic(ev.Err)
		}
		tttc.RedrawAll()
	}
}
