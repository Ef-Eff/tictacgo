package main

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type User struct {
	conn  *websocket.Conn
	data  []Mark
	lobby *Lobby
}

func (u User) lastMark() Mark {
	return u.data[len(u.data) - 1]
}

// Simple json marshalling of a generic message
func (user *User) sendMessage(m Message) {
	msg, _ := json.Marshal(m)

	if err := user.conn.WriteMessage(1, msg); err != nil {
		log.Fatal(err)
	}
}

// Every message will be the players move on the board, this validates it then broadcasts it
func (user *User) readPlay() {
	defer func() {
		user.lobby.disconnect <- user
		user.conn.Close()
	}()
	for {
		_, msg, err := user.conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}

		if user.lobby.users[user] != user.lobby.game.turn {
			user.sendMessage(Message{"error", "Not your turn, shithead."})
			continue
		}

		var m Mark
		if err := json.Unmarshal(msg, &m); err != nil {
			user.sendMessage(Message{"error", "Something went wrong, shithead"})
			log.Fatal(err)
		}

		user.data = append(user.data, m)

		user.lobby.broadcast <- user
	}
}


func Websockets(l *Lobby, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	user := &User{conn: conn, lobby: l}
	l.connect <- user

	go user.readPlay()
}