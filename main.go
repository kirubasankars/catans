package main

import (
	"catans/game"
	"catans/webserver"
)

func main() {
	lobby := game.NewLobby()

	gs := *new(game.GameSetting)
	gs.NumberOfPlayers = 3

	gameId, _ := lobby.CreateGame(gs)
	game := lobby.GetGame(gameId)
	game.Start()

	webserver.StartWebServer()
}
