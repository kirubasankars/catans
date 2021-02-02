package main

import (
	"catans/game"
)

func main() {
	lobby := game.NewLobby()
	lobby.CreateGame(game.GameSetting{NumberOfPlayers: 2})
	//webserver.StartWebServer()
}
