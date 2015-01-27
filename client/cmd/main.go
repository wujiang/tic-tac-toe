package main

import (
	"flag"
	"os/user"

	"github.com/golang/glog"
	"github.com/nsf/termbox-go"
	"github.com/wujiang/tic-tac-toe/client"
	"github.com/wujiang/tic-tac-toe/common"
)

// TODO: send PingMessage to server to keep the connection alive

func main() {
	systemUser, err := user.Current()
	var username string
	if err != nil {
		username = "Unknown"
	} else {
		username = systemUser.Username
	}
	server := flag.String("s", "ws://localhost:8001", "server")
	name := flag.String("u", username, "user name")
	flag.Parse()

	tttc := client.Init(*name)
	defer termbox.Close()

	if err := tttc.Connect(*server); err != nil {
		glog.Fatal("Can not connect to server:")
	}

	tttc.Join()

	go tttc.Listener()

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
				tttc.SendSimpleCMD(ttt.CMD_QUIT)
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				tttc.MoveCursor(ttt.LEFT)
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				tttc.MoveCursor(ttt.DOWN)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				tttc.MoveCursor(ttt.UP)
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				tttc.MoveCursor(ttt.RIGHT)
			case termbox.KeyCtrlJ:
				tttc.Join()
			}

			// vim key bindings
			switch ev.Ch {
			case 'i':
				tttc.PinCursor(ttt.MYRUNE)
			case 'q':
				tttc.SendSimpleCMD(ttt.CMD_QUIT)
				break mainloop
			case 'h':
				tttc.MoveCursor(ttt.LEFT)
			case 'j':
				tttc.MoveCursor(ttt.DOWN)
			case 'k':
				tttc.MoveCursor(ttt.UP)
			case 'l':
				tttc.MoveCursor(ttt.RIGHT)
			case 'O':
				tttc.Join()
			}

		case termbox.EventError:
			glog.Fatal("Termbox EventError:")
		}
		tttc.RedrawAll()
	}
}
