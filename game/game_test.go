package game

import (
	"fmt"
	"testing"
)

func TestGame1(t *testing.T) {

	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 3
	game.UpdateGameSetting(gs)
	game.Start()

	for i := 0; i < gs.NumberOfPlayers*2; i++ {
		locations, _ := game.GetPossibleSettlementLocations()
		game.PutSettlement(locations[1])

		roads, _ := game.GetPossibleRoads()
		game.PutRoad(roads[1])
	}

	for i := 0; i < 5; i++ {
		game.RollDice()
	}

	fmt.Println("")
}
