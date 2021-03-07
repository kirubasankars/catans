package main

import "testing"

func TestGameTrade1(t *testing.T) {
	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
	game.UpdateGameSetting(gs)
	game.Start()
	game.context.Phase = Phase4
	game.context.Players[0].cards = [5]int{2, 4, 5, 5, 6}
	game.context.Players[1].cards = [5]int{4, 2, 6, 7, 8}

	tradeIDs, _ := game.SetupTrade([][2]int{{CardBrick,2}}, [][2]int{{CardWool,1}})

	game.AcceptTrade(1, tradeIDs[0][1])

	game.CompleteTrade(tradeIDs[0][1])

	if game.context.Players[0].cards[CardWool] != 6 {
		t.Log("expected to have 6 wool cards, failed")
		t.Fail()
	}
	if game.context.Players[0].cards[CardBrick] != 2 {
		t.Log("expected to have 2 brick cards, failed")
		t.Fail()
	}

	if game.context.Players[1].cards[CardWool] != 5 {
		t.Log("expected to have 5 wool cards, failed")
		t.Fail()
	}

	if game.context.Players[1].cards[CardBrick] != 4 {
		t.Log("expected to have 5 wool cards, failed")
		t.Fail()
	}
}

func TestGameTradeAcceptAndReject(t *testing.T) {
	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 3
	game.UpdateGameSetting(gs)
	game.Start()
	game.context.Phase = Phase4
	game.context.Players[0].cards = [5]int{2, 4, 5, 5, 6}
	game.context.Players[1].cards = [5]int{4, 2, 6, 7, 8}

	tradeIDs, _ := game.SetupTrade([][2]int{{CardBrick,2}}, [][2]int{{CardWool,1}})

	acceptedTradeID := tradeIDs[0][1]
	rejectedTradeID := tradeIDs[1][1]

	game.AcceptTrade(1, acceptedTradeID)
	game.RejectTrade(2, rejectedTradeID)

	trade := game.context.getTrade(acceptedTradeID)
	if trade.OK != 1 {
		t.Log("expected to be accepted, failed")
		t.Fail()
	}

	trade = game.context.getTrade(rejectedTradeID)
	if trade.OK != -1 {
		t.Log("expected to be rejected, failed")
		t.Fail()
	}
}

func TestGameBankTrade(t *testing.T) {
	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
	game.UpdateGameSetting(gs)
	game.Start()
	game.context.Phase = Phase4
	game.context.Players[0].cards = [5]int{0, 4, 5, 0, 0}

	game.context.Bank.Remove(CardBrick, 4)
	game.context.Bank.Remove(CardWool, 5)

	game.BankTrade([2]int{CardBrick,4}, CardWool)
	if game.context.Players[0].cards[CardWool] != 6 {
		t.Log("expected to have 6 wool cards, failed")
		t.Fail()
	}
	game.context.Players[0].ownPort31 = true
	game.BankTrade([2]int{CardWool,3}, CardLumber)

	if game.context.Players[0].cards[CardWool] != 3 {
		t.Log("expected to have 2 wool cards, failed")
		t.Fail()
	}

	if game.context.Players[0].cards[CardLumber] != 1 {
		t.Log("expected to have 1 lumber card, failed")
		t.Fail()
	}
}