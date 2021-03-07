package main

type GameState struct {
	Phase           string
	Action          GameAction
	Bank            *Bank
	Tiles           [][2]int
	RobberPlacement int

	CurrentPlayerID int
	Players         []*Player

	EventID         int
	Events          []string
}

type Settlement struct {
	Tokens       []int
	TileIndex    []int
	Intersection int
	Upgraded     bool
}
