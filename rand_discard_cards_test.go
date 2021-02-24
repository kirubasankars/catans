package main

import "testing"

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
