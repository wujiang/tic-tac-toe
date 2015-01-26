package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/wujiang/tic-tac-toe/server"
)

func main() {
	addr := flag.String("addr", ":8001", "host address")
	flag.Parse()
	http.HandleFunc("/", server.WSHandler)

	fmt.Println("Server is running at", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
