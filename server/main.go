package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	addr := flag.String("addr", ":8001", "host address")
	flag.Parse()
	http.HandleFunc("/", WSHandler)

	fmt.Println("Server is running at", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		glog.Exitln(err)
	}
}
