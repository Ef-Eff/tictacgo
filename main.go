package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	clargs := ParseCLArgs()
	server := NewServer()
	go server.read()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		RunServer(server, w, r)
	})

	port, envSet := os.LookupEnv("PORT")

	if envSet == false {
		port = strconv.Itoa(*clargs.Port)
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
