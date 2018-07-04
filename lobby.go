package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// The "hub" from the gorilla websocket chat example, ish
type Lobby struct {
	users      map[*User]int
	connect    chan *User
	disconnect chan *User
	broadcast  chan *User
	game       *Game
	server     *Server
	lobbyNum   int
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

func (l *Lobby) writeToAll(m Message) {
	msg, _ := json.Marshal(m)

	for user, _ := range l.users {
		user.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

type Win struct {
	Position     int
	PlayerNumber int
	WinPos       []int
}

func (l *Lobby) endGame(user *User, positions []int) {
	log.Println("Match Finished! Player", l.users[user], "won!")

	res := Win{user.lastPlayedPos(), l.users[user], positions}
	l.writeToAll(Message{"win", res})
}

func (l *Lobby) deleteSelf() {
	l.server.removeLobby <- l.lobbyNum
}

func (l *Lobby) run() {
	defer l.deleteSelf()
	for {
		select {
		case user := <-l.connect:
			l.users[user] = len(l.users) + 1

			user.sendMessage(Message{"welcome", l.users[user]})
			log.Println("A user has connected")

			if len(l.users) == 2 {
				log.Println("Match found!")
				l.newGame()
			}
		case user := <-l.broadcast:
			positions := l.game.play(user)
			if positions != nil {
				l.endGame(user, positions)
				return
			}

			mType := "mark"
			if l.game.counter == 9 {
				mType = "draw"
			}

			res := map[string]int{"Position": user.lastPlayedPos(), "PlayerNumber": l.users[user]}
			l.writeToAll(Message{mType, res})
			if mType == "draw" { return }
		case user := <-l.disconnect:
			delete(l.users, user)

			if l.game != nil {
				log.Println("A user disconnected during a match.")
				l.writeToAll(Message{Type: "winbydc"})
				return
			}
		}
	}
}

func (l *Lobby) shutdown() {
	for user, _ := range l.users {
		user.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "The lobby is shutting down."))
	}
}
