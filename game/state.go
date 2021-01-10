package game

type GameState struct {
	Phase         string
	Action 		  PlayerAction

	CurrentPlayer int
	Players       []*Player
	Roads		  []Road
	Settlements   []Settlement
}

type Settlement struct {
	Tokens       []int
	Indices      []int
	Intersection string

	player       *Player
}

type Road struct {
	Points    [2]string
	player    *Player
}
