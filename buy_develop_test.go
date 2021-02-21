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
	game.context.phase = Phase4
	game.context.Players[0].Cards = [5]int{}

	game.context.Bank.Get(CardWool, 1)
	game.context.Bank.Get(CardGrain, 1)
	game.context.Bank.Get(CardOre, 1)

	game.context.Players[0].Cards[CardWool] = 1
	game.context.Players[0].Cards[CardGrain] = 1
	game.context.Players[0].Cards[CardOre] = 1

	err := game.BuyDevelopmentCard()
	if err != nil {
		t.Log("expected not have error, failed")
		t.FailNow()
	}

	if len(game.getPlayer(game.CurrentPlayer()).DevCards) < 0 {
		t.Log("expected not have dev card, failed")
		t.Fail()
	}

	if game.context.Bank.devCardIndex != 24 {
		t.Log("expected not have dev card, failed")
		t.Fail()
	}
}
