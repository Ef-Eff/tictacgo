package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	defaultPort = 3000
)

// All messages sent from the server are in this generic format
// hopefully using interface{} isnt problematic, i dunno ¯\_(ツ)_/¯
type Message struct {
	Type string
	Data interface{}
}

// This is probably the most useless and shittiest implementation of command line interactivity ever.
// Im just doing it because it's new to me.
// Command Line Arguments (CLArgs®)
type CLArgs struct {
	Port *int
}

func ParseCLArgs() CLArgs {
	port := flag.Int("port", defaultPort, "Set's the port for the server to run on.")
	flag.IntVar(port, "p", defaultPort, "Set's the port for the server to run on.")
	flag.Parse()
	return CLArgs{Port: port}
}

func main() {
	clargs := ParseCLArgs()
	server := newServer()
	go server.read()
	// The next two handle funcs are for the frontend, you can ommit these or somethin if you want to use it as an api
	// Thanks to stackoverflow user RayfenWindspear for below
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		RunServer(server, w, r)
	})

	log.Println("Listening on port", *clargs.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *clargs.Port), nil))
}
