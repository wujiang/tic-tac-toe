package main

import (
	"flag"
	"os/user"

	"github.com/golang/glog"
	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe"
)

func main() {
	systemUser, err := user.Current()
	var username string
	if err != nil {
		username = "Unknown"
	} else {
		username = systemUser.Username
	}
	server := flag.String("s", "ws://localhost:8080", "server")
	name := flag.String("u", username, "user name")
	flag.Parse()

	tttc := TTTCInit(*name)
	defer termbox.Close()

	if err := tttc.Connect(*server); err != nil {
		glog.Exitln("Can not connect to server.")
	}

	go tttc.Listener()

	tttc.RedrawAll()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			// arrows and emacs key bindings
			switch ev.Key {
			case termbox.KeyEnter, termbox.KeySpace:
				tttc.PinCursor(ttt.MyRune)
			case termbox.KeyEsc:
				tttc.Quit()
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				tttc.MoveCursor(ttt.Left)
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				tttc.MoveCursor(ttt.Down)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				tttc.MoveCursor(ttt.Up)
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				tttc.MoveCursor(ttt.Right)
			case termbox.KeyF1:
				tttc.Join(true)
			case termbox.KeyF2:
				tttc.Join(false)
			}

			// vim key bindings
			switch ev.Ch {
			case 'i':
				tttc.PinCursor(ttt.MyRune)
			case 'q':
				tttc.Quit()
				break mainloop
			case 'h':
				tttc.MoveCursor(ttt.Left)
			case 'j':
				tttc.MoveCursor(ttt.Down)
			case 'k':
				tttc.MoveCursor(ttt.Up)
			case 'l':
				tttc.MoveCursor(ttt.Right)
			}

		case termbox.EventError:
			glog.Fatal("Termbox EventError")
		}
		tttc.RedrawAll()
	}
}
