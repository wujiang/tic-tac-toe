package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	addr := flag.String("p", ":8001", "port")
	flag.Parse()
	http.HandleFunc("/", WSHandler)

	fmt.Println("Server is running at", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		glog.Exitln(err)
	}
}
