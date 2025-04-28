package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Lobby struct {
	users      map[*User]int
	connect    chan *User
	disconnect chan *User
	broadcast  chan *User
	game       *Game
	server     *Server
	lobbyNum   int
}

type Win struct {
	Position     int
	PlayerNumber int
	WinPos       []int
}

func (l *Lobby) writeToAll(m Message) {
	msg, _ := json.Marshal(m)

	for user := range l.users {
		user.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (l *Lobby) endGame(user *User) {
	log.Println("Match Finished! Player", l.users[user], "won!")

	res := Win{user.lastPlayedPos(), l.users[user], l.game.winningPos}
	l.writeToAll(Message{WIN, res})
}

func (l *Lobby) deleteSelf() {
	l.server.removeLobby <- l.lobbyNum
}

func (l *Lobby) startNewGame() {
	l.game = NewGame()
	l.writeToAll(Message{Type: START})
}

func (l *Lobby) run() {
	defer l.deleteSelf()
	for {
		select {
		case user := <-l.connect:
			l.users[user] = len(l.users) + 1

			user.sendMessage(Message{WELCOME, l.users[user]})
			log.Println("A user has connected")

			if len(l.users) == 2 {
				log.Println("Match found!")
				l.startNewGame()
			}
		case user := <-l.broadcast:
			l.game.play(user)
			if l.game.winningPos != nil {
				l.endGame(user)
				return
			}

			mType := MARK
			if l.game.counter == 9 {
				mType = DRAW
			}

			res := map[string]int{"Position": user.lastPlayedPos(), "PlayerNumber": l.users[user]}
			l.writeToAll(Message{mType, res})
			if mType == DRAW {
				return
			}
		case user := <-l.disconnect:
			delete(l.users, user)

			if l.game != nil {
				log.Println("A user disconnected during a match.")
				l.writeToAll(Message{Type: WINBYDC})
				return
			}
		}
	}
}

func (l *Lobby) shutdown() {
	for user := range l.users {
		user.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "The lobby is shutting down."))
	}
}
