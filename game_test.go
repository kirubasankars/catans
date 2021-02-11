package main

import (
	"fmt"
	"testing"
)

func TestGameRollDiceAndCards(t *testing.T) {

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
	game.context.handleDice(6)
	game.context.handleDice(8)
	game.context.handleDice(4)

	if game.getPlayer(0).Cards[3] != 1 {
		t.Log("expecting 1 grain, failed")
		t.Fail()
	}
	if game.getPlayer(0).Cards[1] != 1 {
		t.Log("expecting 1 brick, failed")
		t.Fail()
	}
	if game.getPlayer(0).Cards[2] != 1 {
		t.Log("expecting 1 wool, failed")
		t.Fail()
	}

	if game.getPlayer(1).Cards[0] != 1 {
		t.Log("expecting 1 tree, failed")
		t.Fail()
	}
	if game.getPlayer(1).Cards[3] != 2 {
		t.Log("expecting 1 grain, failed")
		t.Fail()
	}

	if game.context.Bank.cards[0] != 18 {
		t.Log("expecting 18 tree remaining, failed")
		t.Fail()
	}

	if game.context.Bank.cards[1] != 18 {
		t.Log("expecting 18 brick remaining, failed")
		t.Fail()
	}

	if game.context.Bank.cards[3] != 16 {
		t.Log("expecting 18 grain remaining, failed")
		t.Fail()
	}

	fmt.Println("")
}

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

func TestGameDiscardCards(t *testing.T) {
	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 3
	gs.DiscardCardLimit = 7
	game.UpdateGameSetting(gs)

	counterFn := func(cards [5]int) int {
		count := 0
		for _, card := range cards {
			count += card
		}
		return count
	}

	game.context.Players[0].Cards = [5]int{3, 3, 4, 5, 6}
	game.context.Players[1].Cards = [5]int{3, 3, 4, 5, 7}
	game.context.Players[2].Cards = [5]int{1, 2, 3, 0, 0}

	game.context.randomDiscardCards()
	n := counterFn(game.context.Players[0].Cards)
	if n != 11 {
		t.Fail()
	}

	n = counterFn(game.context.Players[1].Cards)
	if n != 11 {
		t.Fail()
	}

	n = counterFn(game.context.Players[2].Cards)
	if n != 6 {
		t.Fail()
	}

	game.context.Players[0].Cards = [5]int{3, 3, 4, 5, 25}

	game.context.randomDiscardCards()

	n = counterFn(game.context.Players[0].Cards)
	if n != 20 {
		t.Fail()
	}
}

func TestPlayMonopoly(t *testing.T) {
	game := NewGame()
	game.UpdateGameSetting(GameSetting{NumberOfPlayers: 4, Map: 0})
	game.context.phase = Phase4

	game.context.Players[0].DevCards = append(game.context.Players[0].DevCards, DevCardMonopoly)
	game.context.Players[1].DevCards = append(game.context.Players[1].DevCards, DevCardMonopoly)

	var player0Cards = [5]int{3, 3, 4, 5, 6}
	var player1Cards = [5]int{3, 3, 4, 5, 7}
	var player2Cards = [5]int{1, 2, 3, 0, 0}
	var player3Cards = [5]int{1, 0, 3, 0, 0}

	game.context.Players[0].Cards = player0Cards
	game.context.Players[1].Cards = player1Cards
	game.context.Players[2].Cards = player2Cards
	game.context.Players[3].Cards = player3Cards

	game.context.CurrentPlayerID = 1
	monopolyCardType := 1
	game.context.playMonopoly(monopolyCardType)

	for idx, card := range game.context.Players[0].Cards {
		if idx == monopolyCardType {
			if card != 0 {
				t.Log("expected to have 0 cards, failed")
			}
			continue
		}
		if card != player0Cards[idx] {
			t.Logf("expected to have %d cards, failed", player0Cards[idx])
		}
	}

	for idx, card := range game.context.Players[1].Cards {
		if idx == monopolyCardType {
			if card != 8 {
				t.Log("expected to have 8 cards, failed")
			}
			continue
		}
		if card != player1Cards[idx] {
			t.Logf("expected to have %d cards, failed", player1Cards[idx])
		}
	}

	for idx, card := range game.context.Players[2].Cards {
		if idx == monopolyCardType {
			if card != 0 {
				t.Log("expected to have 0 cards, failed")
			}
			continue
		}
		if card != player2Cards[idx] {
			t.Logf("expected to have %d cards, failed", player2Cards[idx])
		}
	}

	for idx, card := range game.context.Players[3].Cards {
		if idx == monopolyCardType {
			if card != 0 {
				t.Log("expected to have 0 cards, failed")
			}
			continue
		}
		if card != player3Cards[idx] {
			t.Logf("expected to have %d cards, failed", player3Cards[idx])
		}
	}

	if Contains(game.context.Players[1].DevCards, DevCardMonopoly) {
		t.Log("expected monopoly card removed from current player, failed.")
		t.Fail()
	}

	if !Contains(game.context.Players[0].DevCards, DevCardMonopoly) {
		t.Log("expected to have monopoly card, failed.")
		t.Fail()
	}
}

func TestAvailableRoads(t *testing.T) {

}