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

// All messages sent from the server are in this generic format
// hopefully using interface{} isnt problematic, i dunno ¯-_(ツ)_-¯
type Message struct {
	Type string
	Data interface{}
}

// This is the only message format the client is sending (for now)
// I want the array of keys to calculate the score, then the position to disable it
type Mark struct {
	Keys []string `json:"keys"`
	Position int  `json:"position"`
}

type User struct {
	conn *websocket.Conn
	data []Mark
	lobby *Lobby
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
	for {
		_, msg, _ := user.conn.ReadMessage()

		if user.lobby.users[user.lobby.game.turn] != user {
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

func (u User) lastMark() Mark {
	return u.data[len(u.data) - 1]
}

// The "hub" from the gorilla websocket chat example
type Lobby struct {
	users []*User
	connect chan *User
	broadcast chan *User
	game *Game
}

func newLobby() *Lobby {
	return &Lobby{
		users: make([]*User, 0),
		connect: make(chan *User),
		broadcast: make(chan *User),
	}
}

func (l *Lobby) writeToAll(m Message) {
	msg, _ := json.Marshal(m)
	for _, user := range l.users {
		user.conn.WriteMessage(1, msg)
	}
}

func (l *Lobby) run() {
	for {
		select {
		case user := <-l.connect:
			l.users = append(l.users, user)

			user.sendMessage(Message{"welcome", len(l.users)})

			if len(l.users) == 2 {
				log.Println("Match found!")
				l.newGame()
			}
		case user := <- l.broadcast:
			win := l.game.play(user)
			if win == true {
				log.Println("ahoy matey!!!")
				// l.endGame()
				break
			}
			res := map[string]int{"Position": user.lastMark().Position, "PlayerNumber": l.game.turn}
			l.writeToAll(Message{"mark", res})
		}
	}
}

// The Game
type Game struct {
	boardPos map[int]bool 
	turn int
	counter int
	scores map[string]int
}

func (l *Lobby) newGame() {
	game := &Game{
		boardPos: make(map[int]bool),
		turn: 0,
		counter: 0,
		scores: map[string]int{
			"h1": 0, "h2": 0, "h3": 0,
			"v1": 0, "v2": 0, "v3": 0,
			"d1": 0, "d2": 0,
		},
	}
	for i, _ := range game.boardPos {
		game.boardPos[i] = true
	}
	l.writeToAll(Message{Type:"start"})
	l.game = game
}

func (g *Game) play(user *User) bool {
	mark := user.lastMark()
	g.counter++
	for _, v := range mark.Keys {
		if g.turn == 1 { 
			if g.scores[v]++; g.counter > 4 && g.scores[v] == 3 {
				return true
			} 
		} else { 
			if g.scores[v]--; g.counter > 4 && g.scores[v] == -3 {
				return true
			} 
		}
	}
	g.boardPos[mark.Position] = false
	switch g.turn {
	case 0:
		g.turn = 1
	default:
		g.turn = 0
	}
	return false
}

// func (l *Lobby) endGame() {
// 	log.Println("Match Finished! Player", l.game.turn + 1, "won!")
// 	l.writeToAll(Message{"winner", l.users[l.game.turn].lastMark().Position})
// 	poss := make([]int, 0)
// 	for k, v := range l.game.boardPos {
// 		if v == true {
// 			l.game.boardPos[k] = false
// 			poss = append(poss, k)
// 		}
// 	}
// 	js, _ := json.Marshal(Message{"disable", poss})
// 	for _, user := range l.users {
// 		if err := user.conn.WriteMessage(1, js); err != nil {
// 			log.Fatal("In Match:", err)
// 		}
// 	}
// }

// This is probably the most useless and shittiest implementation of command line interactivity ever.
// Im just doing it because it's new to me.
type CLArgs struct {
	Port *int
}

func ParseCLArgs() CLArgs {
	port := flag.Int("port", defaultPort, "Set's the port for the server to run on.")
	flag.IntVar(port, "p", defaultPort, "Set's the port for the server to run on.")
	flag.Parse()
	return CLArgs{Port: port}
}


func Websockets(l *Lobby, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	user := &User{conn: conn, lobby: l}
	l.connect <- user

	go user.readPlay()
}

func main() {
	clargs := ParseCLArgs()
	lobby := newLobby()
	go lobby.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		Websockets(lobby, w, r)
	})

	// Thanks to stackoverflow user RayfenWindspear for below
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Listening on port", *clargs.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", *clargs.Port), nil)
}