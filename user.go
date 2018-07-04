package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	
	"github.com/gorilla/websocket"
)

type User struct {
	conn  *websocket.Conn
	data  []int
	lobby *Lobby
}

func (u User) lastPlayedPos() int {
	return u.data[len(u.data)-1]
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
		_, bytes, err := user.conn.ReadMessage()

		// This error returns when the user has disconnected
		if err != nil {
			log.Println(err)
			break
		}

		if user.lobby.users[user] != user.lobby.game.turn {
			user.sendMessage(Message{"error", "Not your turn, shithead."})
			continue
		}

		if pos, err := strconv.Atoi(string(bytes)); err == nil {
			user.data = append(user.data, pos)

			user.lobby.broadcast <- user
			continue
		}
		
		log.Println(err)
		user.sendMessage(Message{"error", "You aren't sending the right data for some reason, shithead."})
	}
}

func Websockets(l *Lobby, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	user := &User{conn: conn, lobby: l}
	l.connect <- user

	go user.readPlay()
}
