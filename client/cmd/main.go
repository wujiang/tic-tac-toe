package main

import (
	"flag"

	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe/client"
	"github.com/wujiang/tic-tac-toe/common"
)

// TODO: send PingMessage to server to keep the connection alive

func main() {
	server := flag.String("s", "ws://localhost:8001/ws", "server")
	flag.Parse()

	tttc := client.Init()
	defer termbox.Close()

	if err := tttc.Connect(*server); err != nil {
		panic(err)
	}

	tttc.RedrawAll()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			// arrows ane emacs key bindings
			switch ev.Key {
			case termbox.KeyEnter:
				tttc.PinCursor(ttt.MYRUNE)
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				tttc.MoveCursor(ttt.LEFT)
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				tttc.MoveCursor(ttt.DOWN)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				tttc.MoveCursor(ttt.UP)
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				tttc.MoveCursor(ttt.RIGHT)

			}

			// vim key bindings
			switch ev.Ch {
			case 'o':
				tttc.PinCursor(ttt.MYRUNE)
			case 'q':
				break mainloop
			case 'h':
				tttc.MoveCursor(ttt.LEFT)
			case 'j':
				tttc.MoveCursor(ttt.DOWN)
			case 'k':
				tttc.MoveCursor(ttt.UP)
			case 'l':
				tttc.MoveCursor(ttt.RIGHT)
			}

		case termbox.EventError:
			panic(ev.Err)
		}
		tttc.RedrawAll()
	}
}
