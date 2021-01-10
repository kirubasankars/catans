package main

import (
	"catans/game"
	"time"
)

func main() {
	catans := game.NewGameEngine()

	gs := *new(game.GameSetting)
	gs.NumberOfPlayers = 3
	game := catans.CreateGame(gs)

	catans.Start(game)

	time.Sleep(5600 * time.Minute)
}
