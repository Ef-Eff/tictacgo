# TicTacGo

TicTacGo is a Tic Tac Toe/Noughts & Crosses implementation using [Golang](http://golang.org/) with [Gorilla WebSocket](https://github.com/gorilla/websocket) for the backend, and a simple html/css/js frontend. It draws it's structure HEAVILY from the Gorilla WebSocket [chat](https://github.com/gorilla/websocket/tree/master/examples/chat) example.

I made this in an attempt to learn Golang, while creating something I would be interested building upon.

## Actions

* welcome - Give the player their number (1 || 2)
* start - Two players have matched up, the game has started
* mark - BE: The position and player number - FE: Position and keys that entails
* win - A player has won. The position, player number and the key of their win condition (e.g. "d2")
* winbydc - The match has started and a player has disconnected.
* error - Something went wrong
