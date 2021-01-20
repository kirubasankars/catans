package main

import (
	"catans/game"
	"catans/webserver"
)

func main() {
	lobby := game.NewLobby()

	gs := *new(game.GameSetting)
	gs.Map = "default"
	gs.NumberOfPlayers = 3

	gameId := lobby.CreateGame(gs)
	game := lobby.GetGame(gameId)
	game.Start()

	webserver.StartWebServer()
}
