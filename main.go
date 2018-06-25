package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"flag"
)

// A bit of a useless struct, just need to check that the template works
type Page struct {
	Title string
	Body string
}

// This is probably the most useless and shittiest implementation of command line interactivity ever.
// Im just doing it because it's new to me.
type CLArgs struct {
	Port *int
}

func ParseCLArgs() CLArgs {
	port := flag.Int("port", 3000, "Set's the port for the server to run on.")
	flag.IntVar(port, "p", 3000, "Set's the port for the server to run on.")
	flag.Parse()
	return CLArgs{Port: port}
}

// The web page the game is played on
func Index(w http.ResponseWriter, r *http.Request) {
	data := Page{"Hello, World!", "Lorem ipsum something something idk"}
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, data)
}

func main() {
	clargs := ParseCLArgs()
	defer http.ListenAndServe(fmt.Sprintf(":%v", *clargs.Port), nil)
	defer log.Println("Listening on port", *clargs.Port)

	http.HandleFunc("/", Index)
}