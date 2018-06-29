package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	defaultPort = 3000
)

type Server struct {
	lobbies []*Lobby
}

func newServer() *Server {
	return &Server{make([]*Lobby,0)}
}

// All messages sent from the server are in this generic format
// hopefully using interface{} isnt problematic, i dunno ¯\_(ツ)_/¯
type Message struct {
	Type string
	Data interface{}
}

// This is the only message format the client is sending (for now)
// I want the array of keys to calculate the score, then the position to disable it
type Mark struct {
	Keys 		 []string `json:"keys"`
	Position int  		`json:"position"`
}

// This is probably the most useless and shittiest implementation of command line interactivity ever.
// Im just doing it because it's new to me.
type CLArgs struct {
	Port *int
}

func ParseCLArgs() CLArgs {
	port := flag.Int("port", defaultPort, "Set's the port for the server to run on.")
	flag.IntVar(port, "p", defaultPort, "Set's the port for the server to run on.")
	flag.Parse()
	return CLArgs{Port: port}
}

func main() {
	clargs := ParseCLArgs()
	server := newServer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var lobby *Lobby
		if len(server.lobbies) == 0 || len(server.lobbies[len(server.lobbies)-1].users) == 2 {
			lobby = newLobby()
			server.lobbies = append(server.lobbies, lobby)
			go lobby.run()
		} else {
			lobby = server.lobbies[len(server.lobbies)-1]
		}
		Websockets(lobby, w, r)
	})

	// Thanks to stackoverflow user RayfenWindspear for below
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Listening on port", *clargs.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *clargs.Port), nil))
}