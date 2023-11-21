package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	WINBYDC MessageType = "winbydc"
	DRAW    MessageType = "draw"
	MARK    MessageType = "mark"
	WELCOME MessageType = "welcome"
	WIN     MessageType = "win"
	ERROR   MessageType = "error"
	START   MessageType = "start"
)

type MessageType string

// All messages sent from the server are in this generic format
type Message struct {
	Type MessageType
	Data interface{}
}

func main() {
	clargs := ParseCLArgs()
	server := NewServer()
	go server.read()

	// Thanks to stackoverflow user RayfenWindspear for bellow
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		RunServer(server, w, r)
	})

	log.Println("Listening on port", *clargs.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *clargs.Port), nil))
}
