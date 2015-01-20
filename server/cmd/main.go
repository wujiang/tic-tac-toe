package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	addr := flag.String("addr", ":8001", "host address")
	flag.Parse()

	http.HandleFunc("/ws", server.WSHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		glog.Fatal("Can not start server:", err)
	}
}
