package main

import (
	"catans/board"
	"catans/webserver"
)

func main() {
	board.NewBoard("default")
	//engine := game.NewGameEngine()

	//gs := *new(game.GameSetting)
	//gs.Map = "default"
	//gs.NumberOfPlayers = 3
	//
	//gameId := engine.CreateGame(gs)
	//engine.Start(gameId)

	webserver.StartWebServer()
}
