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

func (l *Lobby) endGame(user *User, positions []int) {
	log.Println("Match Finished! Player", l.users[user], "won!")

	res := Win{user.lastPlayedPos(), l.users[user], positions}
	l.writeToAll(Message{WIN, res})
}

func (l *Lobby) deleteSelf() {
	l.server.removeLobby <- l.lobbyNum
}

func (l *Lobby) newGame() {
	game := &Game{
		boardPos: make(map[int]bool, 9),
		turn:     1,
		counter:  0,
		scores: map[string]int{
			"h1": 0, "h2": 0, "h3": 0,
			"v1": 0, "v2": 0, "v3": 0,
			"d1": 0, "d2": 0,
		},
	}

	for i := range game.boardPos {
		game.boardPos[i] = true
	}

	l.writeToAll(Message{Type: START})

	l.game = game
}

// Split into seperate functions for each action for clarity?
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
				l.newGame()
			}
		case user := <-l.broadcast:
			positions := l.game.play(user)
			if positions != nil {
				l.endGame(user, positions)
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
