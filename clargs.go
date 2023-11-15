package main

import "flag"

const (
	defaultPort = 3000
	portMessage = "Set the port for the server to run on."
)

type CLArgs struct {
	Port *int
}

func parseCLArgs() CLArgs {
	port := flag.Int("port", defaultPort, portMessage)
	flag.IntVar(port, "p", defaultPort, portMessage)
	flag.Parse()
	return CLArgs{Port: port}
}
