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
	boardPos map[int]bool 
	turn int
	counter int
	scores map[string]int
}

func newGame(l *Lobby) *Game {
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
	return game
}

func (g *Game) play(user *User) bool {
	mark := user.lastMark()
	g.counter++
	for _, v := range mark.Score {
		if g.turn == 1 { 
			if g.scores[v]++; g.scores[v] == 3 {
				return true
			} 
		} else { 
			if g.scores[v]--; g.scores[v] == -3 {
				return true
			} 
		}
	}
	log.Println(g.scores)
	g.boardPos[mark.Position] = false
	switch g.turn {
	case 0:
		g.turn = 1
	default:
		g.turn = 0
	}
	return false
}

type ErrResponse struct {
	Type string
	Message string
}

type RegularResponse struct {
	Type string
	Position int
}

type DisableResponse struct {
	Type string
	Positions []int
}

type Mark struct {
	Score []string `json:"score"`
	Position int   `json:"position"`
}

type User struct {
	conn *websocket.Conn
	data []Mark
	lobby *Lobby
}

func (u User) lastMark() Mark {
	return u.data[len(u.data) - 1]
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

type Player struct {
	Player int
}

func (l *Lobby) writeToAll(r RegularResponse) {
	js, _ := json.Marshal(r)
	for _, user := range l.match {
		if err := user.conn.WriteMessage(1, js); err != nil {
			log.Fatal("In Match:", err)
		}
	}
}

func (l *Lobby) run() {
	for {
		select {
		case user := <-l.connect:
			l.match = append(l.match, user)
			log.Println("User Connected")
			js, _ :=  json.Marshal(Player{len(l.match)})
			if err := user.conn.WriteMessage(1, js); err != nil {
				log.Fatal("????:", err)
			}
			if len(l.match) == 2 {
				log.Println("Match found!")
				l.game = newGame(l)
			}
		case user := <- l.broadcast:
			win := l.game.play(user)
			if win != false {
				l.endGame()
				break
			}
			l.writeToAll(RegularResponse{"porra", user.lastMark().Position})
		}
	}
}

func (l *Lobby) endGame() {
	log.Println("Match Finished! Player", l.game.turn + 1, "won!")
	l.writeToAll(RegularResponse{"winner", l.match[l.game.turn].lastMark().Position})
	poss := make([]int, 0)
	for k, v := range l.game.boardPos {
		if v == true {
			l.game.boardPos[k] = false
			poss = append(poss, k)
		}
	}
	js, _ := json.Marshal(DisableResponse{"disable", poss})
	for _, user := range l.match {
		if err := user.conn.WriteMessage(1, js); err != nil {
			log.Fatal("In Match:", err)
		}
	}
}

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

	go func(user *User) {
		for {
			_, msg, _ := user.conn.ReadMessage()
			if l.match[l.game.turn] != user {
				jss, _ := json.Marshal(ErrResponse{"error", "Not your turn shithead"})
				user.conn.WriteMessage(1, jss)
				continue
			}
			var m Mark
			if err := json.Unmarshal(msg, &m); err != nil {
				jss, _ := json.Marshal(ErrResponse{"error", "Something went wrong, shithead"})
				user.conn.WriteMessage(1, jss)
				log.Fatal(err)
				continue
			}
			user.data = append(user.data, m)
			l.broadcast <- user
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