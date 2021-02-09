package main

func main() {
	lobby := NewLobby()
	lobby.CreateGame(Setting{NumberOfPlayers: 2})
	//webserver.StartWebServer()
}
