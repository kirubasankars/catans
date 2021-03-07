package main

import (
	"testing"
)

func TestBuyDevelopmentCard(t *testing.T) {
	game := NewGame()

	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
	game.UpdateGameSetting(gs)
	game.Start()
	game.context.Phase = Phase4
	game.context.Players[0].cards = [5]int{}

	game.context.Bank.Remove(CardWool, 1)
	game.context.Bank.Remove(CardGrain, 1)
	game.context.Bank.Remove(CardOre, 1)

	game.context.Players[0].cards[CardWool] = 1
	game.context.Players[0].cards[CardGrain] = 1
	game.context.Players[0].cards[CardOre] = 1

	counterFn := func(cards [5]int) int {
		count := 0
		for _, card := range cards {
			count += card
		}
		return count
	}

	err := game.BuyDevelopmentCard()
	if err != nil {
		t.Log("expected not have error, failed")
		t.FailNow()
	}

	if len(game.getPlayer(game.CurrentPlayer()).devCards) <= 0 {
		t.Log("expected to have dev card, failed")
		t.Fail()
	}

	if game.context.Bank.devCardIndex != 24 {
		t.Log("expected not have dev card, failed")
		t.Fail()
	}

	if counterFn(game.context.Players[0].cards) != 0 {
		t.Log("expected to have card removed, failed")
		t.Fail()
	}
}

func TestBuyDevelopmentCardError(t *testing.T) {
	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
	game.UpdateGameSetting(gs)
	game.Start()
	game.context.Phase = Phase4
	game.context.Players[0].cards = [5]int{}

	game.context.Bank.Remove(CardWool, 1)
	game.context.Bank.Remove(CardGrain, 1)
	game.context.Bank.Remove(CardOre, 1)

	err := game.BuyDevelopmentCard()
	if err == nil {
		t.Log("expected to have err, failed")
		t.Fail()
	}
}
