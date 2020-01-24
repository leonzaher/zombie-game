package main

import (
	"flag"
	"fmt"
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
	log.Println(fmt.Sprintf("Starting the server at port %s", *addr))

	err := http.ListenAndServe(":"+*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
