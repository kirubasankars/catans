package game

import (
	"fmt"
	"testing"
)

func TestGame1(t *testing.T) {

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

	player0 := game.getPlayer(0)
	if len(player0.Settlements) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}
	if len(player0.Roads) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}

	player1 := game.getPlayer(1)
	if len(player1.Settlements) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}
	if len(player1.Roads) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}

	if player0.Settlements[0].Intersection != 14 {
		t.Log("expecting settlement in 14, failed")
		t.Fail()
	}
	if player0.Settlements[1].Intersection != 13 {
		t.Log("expecting settlement in 13, failed")
		t.Fail()
	}

	if player1.Settlements[0].Intersection != 26 {
		t.Log("expecting settlement in 14, failed")
		t.Fail()
	}
	if player1.Settlements[1].Intersection != 41 {
		t.Log("expecting settlement in 13, failed")
		t.Fail()
	}

	game.context.handleDice(12)

	if game.getPlayer(0).Cards[3] != 1 {
		t.Log("expecting 1 grain, failed")
		t.Fail()
	}

	if game.context.Bank.cards[3] != 18 {
		t.Log("expecting 18 grain remaining, failed")
		t.Fail()
	}

	game.context.handleDice(6)

	if game.getPlayer(0).Cards[1] != 1 {
		t.Log("expecting 1 grain, failed")
		t.Fail()
	}

	if game.context.Bank.cards[1] != 18 {
		t.Log("expecting 18 grain remaining, failed")
		t.Fail()
	}

	if game.getPlayer(1).Cards[3] != 1 {
		t.Log("expecting 1 grain, failed")
		t.Fail()
	}

	if game.context.Bank.cards[3] != 17 {
		t.Log("expecting 18 grain remaining, failed")
		t.Fail()
	}

	fmt.Println("")
}

func TestGame2(t *testing.T) {

	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
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

	game.context.handleDice(7)

	fmt.Println("")
}