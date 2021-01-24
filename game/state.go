package game

type GameState struct {
	Phase         string
	Action        gameAction
	Bank          *Bank
	CurrentPlayer int
	Players       []*Player
}

type Settlement struct {
	Tokens       []int
	Indices      []int
	Intersection int
}

type Road struct {
	Points    [2]int
	player    *Player
}


