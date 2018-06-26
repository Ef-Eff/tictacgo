package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	defaultPort = 3000
)

type Game struct {
	playerOne *User
	playerTwo *User
	board []Mark
	turn int
}

func newGame(l *Lobby) *Game {
	return &Game{
		playerOne: l.match[0],
		playerTwo: l.match[1],
		board: make([]Mark, 0),
		turn: 0,
	}
}

type ErrResponse struct {
	Type string
	Message string
}

type Mark struct {
	Type string
	Player int
	Position string
}

func (g *Game) play(user *User) {
	result := Mark{ "play", g.turn, user.lastMessage() }
	g.board = append(g.board, result)
	log.Println(g.board)
	switch g.turn {
	case 0:
		g.turn = 1
	default:
		g.turn = 0
	}
}

type Lobby struct {
	// Will only have a single match for now, which consists of two users
	match []*User
	connect chan *User
	broadcast chan *User
	game *Game
}

func newLobby() *Lobby {
	return &Lobby{
		match: make([]*User, 0),
		connect: make(chan *User),
		broadcast: make(chan *User),
	}
}

func (l *Lobby) run() {
	for {
		select {
		case user := <-l.connect:
			l.match = append(l.match, user)
			log.Println("User Connected")
			if len(l.match) == 2 {
				log.Println("Match found!")
				l.game = newGame(l)
			}
		case user := <- l.broadcast:
			if l.game != nil {
				l.game.play(user)
				js, _ := json.Marshal(l.game.board)
				for _, user := range l.match {
					if err := user.conn.WriteMessage(1, js); err != nil {
						log.Fatal("In Match:", err)
					}
				}
			} else {
				log.Println("Incoming:", user.lastMessage())
				for _, user := range l.match {
					err := user.conn.WriteMessage(1, []byte(user.lastMessage()))
					if err != nil {
						log.Fatal("Broadcast:", err)
					}
				}
			}
		}
	}
}

// This is probably the most useless and shittiest implementation of command line interactivity ever.
// Im just doing it because it's new to me.
type CLArgs struct {
	Port *int
}

type User struct {
	conn *websocket.Conn
	data []string
	lobby *Lobby
}

func (u User) lastMessage() string {
	return u.data[len(u.data) - 1]
}

func ParseCLArgs() CLArgs {
	port := flag.Int("port", defaultPort, "Set's the port for the server to run on.")
	flag.IntVar(port, "p", defaultPort, "Set's the port for the server to run on.")
	flag.Parse()
	return CLArgs{Port: port}
}

func Websockets(l *Lobby, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	user := &User{conn, make([]string, 0), l}
	l.connect <- user

	go func(user *User) {
		for {
			_, msg, _ := user.conn.ReadMessage()
			if l.game != nil && l.match[l.game.turn] != user {
				jss, _ := json.Marshal(ErrResponse{"error", "Not your turn shithead"})
				user.conn.WriteMessage(1, jss)
				continue
			}
			user.data = append(user.data, string(msg))
			l.broadcast <- user
			log.Println(user.data)
		}
	}(user)
}

func main() {
	clargs := ParseCLArgs()
	defer http.ListenAndServe(fmt.Sprintf(":%v", *clargs.Port), nil)
	defer log.Println("Listening on port", *clargs.Port)
	lobby := newLobby()
	go lobby.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		Websockets(lobby, w, r)
	})
}