package main

import (
	"catans/game"
)

func main() {
	lobby := game.NewLobby()
	lobby.CreateGame(game.Setting{NumberOfPlayers: 2})
	//webserver.StartWebServer()
}
