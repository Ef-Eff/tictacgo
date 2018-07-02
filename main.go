package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	mutex    = sync.Mutex{}
)

const (
	defaultPort = 3000
)

type Server struct {
	lobbies     map[int]*Lobby
	lastLobby   int
	removeLobby chan int
}

func newServer() *Server {
	return &Server{
		lobbies:     make(map[int]*Lobby),
		removeLobby: make(chan int),
	}
}

func (s *Server) read() {
	for {
		select {
		case lobbyNum := <-s.removeLobby:
			log.Println("Removing", lobbyNum)
			delete(s.lobbies, lobbyNum)
		}
	}
}

// All messages sent from the server are in this generic format
// hopefully using interface{} isnt problematic, i dunno ¯\_(ツ)_/¯
type Message struct {
	Type string
	Data interface{}
}

func RunServer(server *Server, w http.ResponseWriter, r *http.Request) {
	var lobby *Lobby
	if server.lobbies[server.lastLobby] == nil || len(server.lobbies[server.lastLobby].users) == 2 {
		// Im not sure this is legit use of a mutex.
		// I basically just want to make sure the number is being incremented and the lobby is being created one at a time
		// Hopefully this works. I doubt it would be easy for a collision to happen anyway.
		mutex.Lock()
		server.lastLobby++
		lobby = newLobby(server)
		mutex.Unlock()

		server.lobbies[server.lastLobby] = lobby
		log.Println("Starting lobby", lobby.lobbyNum)
		go lobby.run()
	} else {
		lobby = server.lobbies[server.lastLobby]
	}
	Websockets(lobby, w, r)
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
	go server.read()
	// The next two handle funcs are for the frontend, you can ommit these or somethin if you want to use it as an api
	// Thanks to stackoverflow user RayfenWindspear for below
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
