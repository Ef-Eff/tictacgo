package main

import (
	"encoding/json"
	"log"
)

// The "hub" from the gorilla websocket chat example, ish
type Lobby struct {
	users 		 map[*User]int
	connect 	 chan *User
	disconnect chan *User
	broadcast  chan *User
	game 			 *Game
	server 		 *Server
	lobbyNum   int           
}

func newLobby(s *Server) *Lobby {
	return &Lobby{
		users: 		  make(map[*User]int),
		connect: 	  make(chan *User),
		broadcast:  make(chan *User),
		disconnect: make(chan *User),
		server: s,
		lobbyNum: s.lastLobby,
	}
}

func (l *Lobby) writeToAll(m Message) {
	msg, _ := json.Marshal(m)
	
	for user, _ := range l.users {
		user.conn.WriteMessage(1, msg)
	}
}

type Win struct {
	Position int
	PlayerNumber int
	Key string
}

func (l *Lobby) endGame(user *User, key string) {
	log.Println("Match Finished! Player", l.users[user], "won!")

	res := Win{user.lastMark().Position, l.users[user], key}
	l.writeToAll(Message{"win", res})
	// Self destruct sequence activated
	l.deleteSelf()
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

			if len(l.users) == 2 {
				log.Println("Match found!")
				l.newGame()
			}
		case user := <- l.broadcast:
			key := l.game.play(user)
			if key != "" {
				l.endGame(user, key)
				break
			}
			res := map[string]int{"Position": user.lastMark().Position, "PlayerNumber": l.users[user]}
			l.writeToAll(Message{"mark", res})
		case user := <- l.disconnect:
			delete(l.users, user)
			if l.game != nil {
				log.Println("meme")
				l.writeToAll(Message{Type: "winbydc"})
				return
			}
		}
	}
}
