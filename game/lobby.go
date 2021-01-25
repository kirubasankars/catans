package game

import "fmt"

type GameEngine struct {
	games       map[string]Game
	gameCounter int
}

func (ge GameEngine) CreateGame(gs GameSetting) (string, error) {
	game := NewGame()
	id := fmt.Sprintf("GAME#%d", ge.gameCounter)
	ge.games[id] = *game
	if err := game.UpdateGameSetting(gs); err != nil {
		return "", err
	}
	return id, nil
}

func (ge GameEngine) GetGame(gameID string) Game {
	return ge.games[gameID]
}

func NewLobby() *GameEngine {
	ge := new(GameEngine)
	ge.games = make(map[string]Game)
	ge.gameCounter = 1
	return ge
}
