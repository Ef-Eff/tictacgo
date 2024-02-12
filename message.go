package main

type MessageType string

type Message struct {
	Type MessageType
	Data interface{}
}

const (
	WINBYDC MessageType = "winbydc"
	DRAW    MessageType = "draw"
	MARK    MessageType = "mark"
	WELCOME MessageType = "welcome"
	WIN     MessageType = "win"
	ERROR   MessageType = "error"
	START   MessageType = "start"
)
