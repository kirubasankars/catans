package main

import "fmt"

type GameLobby struct {
	games       map[string]Game
	gameCounter int
}

func (gameLobby *GameLobby) CreateGame(gs GameSetting) (string, error) {
	game := NewGame()
	id := fmt.Sprintf("GAME-%d", gameLobby.gameCounter)
	gameLobby.games[id] = *game
	gameLobby.gameCounter++
	if err := game.UpdateGameSetting(gs); err != nil {
		return "", err
	}
	return id, nil
}

func (gameLobby GameLobby) GetGame(gameID string) Game {
	return gameLobby.games[gameID]
}

func NewLobby() *GameLobby {
	gameLobby := new(GameLobby)
	gameLobby.games = make(map[string]Game)
	gameLobby.gameCounter = 1
	return gameLobby
}
