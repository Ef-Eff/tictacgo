package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	defaultPort = 3000
)

// This is probably the most useless and shittiest implementation of command line interactivity ever.
// Im just doing it because it's new to me.
type CLArgs struct {
	Port *int
}

func newLobby() *Lobby {
	return &Lobby{
		match: make([]*User, 2),
		connect: make(chan *User),
		broadcast: make(chan []byte),
	}
}

type Lobby struct {
	// Will only have a single match for now, which consists of two users
	match []*User
	connect chan *User
	broadcast chan []byte
}

func (l *Lobby) addUser(u *User) {
	l.connect <- u
	l.match = append(l.match, u)
}

func (l *Lobby) run() {
	for {
		select {
		case user := <-l.connect:
			l.match = append(l.match, user)
		case <- l.broadcast:
			for _, user := range l.match {
				err := user.conn.WriteMessage(2, []byte("Whaaaat!"))
				if err != nil {
					user.conn.Close()
					log.Fatal("Broadcast:", err)
				}
			}
		}
	}
}

type User struct {
	conn *websocket.Conn
	data []string
	lobby *Lobby
}

func ParseCLArgs() CLArgs {
	port := flag.Int("port", defaultPort, "Set's the port for the server to run on.")
	flag.IntVar(port, "p", defaultPort, "Set's the port for the server to run on.")
	flag.Parse()
	return CLArgs{Port: port}
}

// The web page the game is played on
func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func Websockets(l *Lobby, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	user := &User{conn, make([]string, 0), l}
	l.addUser(user)
	log.Println("User Connected")

	go func(user *User) {
		for {
			_, msg, _ := user.conn.ReadMessage()
			l.broadcast <- msg
			user.data = append(user.data, string(msg))
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

	http.HandleFunc("/", Index)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		Websockets(lobby, w, r)
	})
}