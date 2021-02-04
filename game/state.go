package game

type State struct {
	Phase           string
	Action          gameAction
	Bank            *Bank
	CurrentPlayerID int
	Players         []*Player
	Tiles           [][2]int
	RobberPlacement int
}

type Settlement struct {
	Tokens       []int
	Indices      []int
	Intersection int
	Upgraded     bool
}

type Road struct {
	Points [2]int
	player *Player
}
