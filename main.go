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
			// arrows ane emacs key bindings
			switch ev.Key {
			case termbox.KeyEnter:
				tttc.PinCursor(client.MYRUNE)
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				tttc.MoveCursor(client.LEFT)
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				tttc.MoveCursor(client.DOWN)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				tttc.MoveCursor(client.UP)
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				tttc.MoveCursor(client.RIGHT)

			}

			// vim key bindings
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
