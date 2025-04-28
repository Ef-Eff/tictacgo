package main

var boardPosKeys = map[int][]string{
	0: {"h1", "v1", "d1"}, 1: {"h1", "v2"}, 2: {"h1", "v3", "d2"},
	3: {"h2", "v1"}, 4: {"h2", "v2", "d1", "d2"}, 5: {"h2", "v3"},
	6: {"h3", "v1", "d2"}, 7: {"h3", "v2"}, 8: {"h3", "v3", "d1"},
}

var winConditions = map[string][]int{
	"h1": {0, 1, 2}, "h2": {3, 4, 5}, "h3": {6, 7, 8},
	"v1": {0, 3, 6}, "v2": {1, 4, 7}, "v3": {2, 5, 8},
	"d1": {0, 4, 8}, "d2": {2, 4, 6},
}

// The Game
type Game struct {
	boardPos   map[int]bool
	turn       int
	counter    int
	winningPos []int
	scores     map[string]int
}

func NewGame() *Game {
	return &Game{
		boardPos: map[int]bool{
			0: true, 1: true, 2: true,
			3: true, 4: true, 5: true,
			6: true, 7: true, 8: true,
		},
		turn:       1,
		counter:    0,
		winningPos: nil,
		scores: map[string]int{
			"h1": 0, "h2": 0, "h3": 0,
			"v1": 0, "v2": 0, "v3": 0,
			"d1": 0, "d2": 0,
		},
	}
}

func (g *Game) changeScores(pos, diff int) {
	for _, key := range boardPosKeys[pos] {
		if g.scores[key] += diff; g.counter > 4 && g.scores[key] == 3*diff {
			g.winningPos = winConditions[key]
		}
	}
}

func (g *Game) flipTurn() int {
	switch g.turn {
	case 1:
		g.turn = 2
		return 1
	default:
		g.turn = 1
		return -1
	}
}

func (g *Game) play(user *User) {
	pos := user.lastPlayedPos()

	g.boardPos[pos] = false
	g.counter++
	diff := g.flipTurn()

	g.changeScores(pos, diff)
}
