package game

import "fmt"

type GameEngine struct {
	games 		map[string]*Game
	gameCounter int
}

func (ge GameEngine) CreateGame(gs GameSetting) string {
	game := NewGame()
	id := fmt.Sprintf("GAME#%d", ge.gameCounter)
	ge.games[id] = game
	game.UpdateGameSetting(gs)
	return id
}

func (ge GameEngine) Start(gameID string) {
	ge.games[gameID].Start()
}

func NewGameEngine() *GameEngine {
	ge := new(GameEngine)
	ge.games = make(map[string]*Game)
	ge.gameCounter = 1
	return ge
}
