package main

import "testing"

func TestPlayMonopoly(t *testing.T) {
	game := NewGame()
	game.UpdateGameSetting(GameSetting{NumberOfPlayers: 4, Map: 0})
	game.context.Phase = Phase4

	game.context.Players[0].devCards = append(game.context.Players[0].devCards, DevCardMonopoly)
	game.context.Players[1].devCards = append(game.context.Players[1].devCards, DevCardMonopoly)
	game.context.Players[1].devCards = append(game.context.Players[0].devCards, DevCardKnight)

	var player0Cards = [5]int{3, 3, 4, 5, 6}
	var player1Cards = [5]int{3, 3, 4, 5, 7}
	var player2Cards = [5]int{1, 2, 3, 0, 0}
	var player3Cards = [5]int{1, 0, 3, 0, 0}

	game.context.Players[0].cards = player0Cards
	game.context.Players[1].cards = player1Cards
	game.context.Players[2].cards = player2Cards
	game.context.Players[3].cards = player3Cards

	game.context.CurrentPlayerID = 1
	monopolyCardType := CardBrick
	game.context.playMonopoly(monopolyCardType)

	for idx, card := range game.context.Players[0].cards {
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

	for idx, card := range game.context.Players[1].cards {
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

	for idx, card := range game.context.Players[2].cards {
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

	for idx, card := range game.context.Players[3].cards {
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

	if Contains(game.context.Players[1].devCards, DevCardMonopoly) {
		t.Log("expected monopoly card removed from current player, failed.")
		t.Fail()
	}

	if !Contains(game.context.Players[0].devCards, DevCardMonopoly) {
		t.Log("expected to have monopoly card, failed.")
		t.Fail()
	}
}
