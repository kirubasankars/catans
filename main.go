package main

func main() {
	lobby := NewLobby()
	lobby.CreateGame(GameSetting{NumberOfPlayers: 2})
	StartWebServer()
}
