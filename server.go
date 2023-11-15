package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	mutex    = sync.Mutex{}
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
			log.Println("Removing lobby:", lobbyNum)
			s.lobbies[lobbyNum].shutdown()
			delete(s.lobbies, lobbyNum)
		}
	}
}

func newLobby(s *Server) *Lobby {
	return &Lobby{
		users:      make(map[*User]int, 2),
		connect:    make(chan *User, 1),
		broadcast:  make(chan *User, 1),
		disconnect: make(chan *User, 1),
		server:     s,
		lobbyNum:   s.lastLobby,
	}
}

func RunServer(server *Server, w http.ResponseWriter, r *http.Request) {
	var lobby *Lobby
	if server.lobbies[server.lastLobby] == nil || len(server.lobbies[server.lastLobby].users) == 2 {
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
