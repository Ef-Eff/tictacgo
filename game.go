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
	boardPos map[int]bool
	turn     int
	counter  int
	scores   map[string]int
}

func (g *Game) changeScores(pos, diff int) (positions []int) {
	for _, key := range boardPosKeys[pos] {
		if g.scores[key] += diff; g.counter > 4 && g.scores[key] == 3*diff {
			positions = winConditions[key]
		}
	}
	return
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

func (g *Game) play(user *User) []int {
	g.counter++

	diff := g.flipTurn()

	pos := user.lastPlayedPos()
	g.boardPos[pos] = false

	return g.changeScores(pos, diff)
}
