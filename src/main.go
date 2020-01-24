package main

import (
	"flag"
	"log"
	"net/http"
	"zombie-game/src/ws"
)

var addr = flag.String("addr", "8080", "http service address")

func main() {
	flag.Parse()
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(":"+*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
