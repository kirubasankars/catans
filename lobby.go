package main

type GameLobby struct {
	games       map[int]Game
	users 		map[string]User
	gameCounter int
}

func (gameLobby *GameLobby) CreateGame(gs GameSetting) (int, error) {
	game := NewGame()
	id := gameLobby.gameCounter
	gameLobby.games[id] = *game
	gameLobby.gameCounter++
	if err := game.UpdateGameSetting(gs); err != nil {
		return -1, err
	}
	return id, nil
}

func (gameLobby GameLobby) GetGame(gameID int) Game {
	return gameLobby.games[gameID]
}

func (gameLobby *GameLobby) AddUserToGame(gameID int, userID string) {
	game := gameLobby.games[gameID]
	user := gameLobby.users[userID]
	game.context.Users = append(game.context.Users, &user)
}

func NewLobby() *GameLobby {
	gameLobby := new(GameLobby)
	gameLobby.games = make(map[int]Game)
	gameLobby.users = make(map[string]User)
	gameLobby.gameCounter = 1
	return gameLobby
}
