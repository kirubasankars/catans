package main

import "testing"

func TestGameFreeCards(t *testing.T) {

	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
	gs.TurnTimeOut = false
	game.UpdateGameSetting(gs)
	game.Start()

	//player 0
	game.PutSettlement(14)
	game.PutRoad([2]int{14, 15})
	//player 1
	game.PutSettlement(26)
	game.PutRoad([2]int{25, 26})

	//player 1
	game.PutSettlement(41)
	game.PutRoad([2]int{41, 32})
	//player 0
	game.PutSettlement(13)
	game.PutRoad([2]int{13, 20})

}