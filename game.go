package main

// The Game
type Game struct {
	boardPos map[int]bool 
	turn 		 int
	counter  int
	scores 	 map[string]int
}

func (l *Lobby) newGame() {
	game := &Game{
		boardPos: make(map[int]bool),
		turn: 		1,
		counter: 	0,
		scores: 	map[string]int{
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

func (g *Game) flipTurn() {
	switch g.turn {
	case 1:
		g.turn = 2
	default:
		g.turn = 1
	}
}

func (g *Game) play(user *User) string {
	mark := user.lastMark()

	g.counter++

	for _, v := range mark.Keys {
		if g.turn == 1 { 
			if g.scores[v]++; g.counter > 4 && g.scores[v] == 3 {
				return v
			} 
		} else { 
			if g.scores[v]--; g.counter > 4 && g.scores[v] == -3 {
				return v
			} 
		}
	}

	g.flipTurn()
	g.boardPos[mark.Position] = false

	return ""
}
