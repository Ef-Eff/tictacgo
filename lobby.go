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
}

func newLobby() *Lobby {
	return &Lobby{
		users: 		 make(map[*User]int),
		connect: 	 make(chan *User),
		broadcast: make(chan *User),
	}
}

func (l *Lobby) writeToAll(m Message) {
	msg, _ := json.Marshal(m)
	
	for user, _ := range l.users {
		user.conn.WriteMessage(1, msg)
	}
}

func (l *Lobby) run() {
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
		case <- l.disconnect:
			l.writeToAll(Message{"error", "someone dced pce"})
		}
	}
}
