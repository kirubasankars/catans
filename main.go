package main

import (
	"catans/game"
	"time"
)

func main() {
	engine := game.NewGameEngine()

	gs := *new(game.GameSetting)
	gs.NumberOfPlayers = 3
	game := engine.CreateGame(gs)
	engine.Start(game)

	time.Sleep(60 * time.Minute)
}
