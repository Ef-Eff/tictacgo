package main

import (
	"sync"
	"net/http"
	"log"
	
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
			log.Println("Removing lobby", lobbyNum)
			s.lobbies[lobbyNum].shutdown()
			delete(s.lobbies, lobbyNum)
		}
	}
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