package main

import "errors"

type GameTrade struct {
	ID          int
	From        int
	To          int
	Gives       [][2]int
	Wants       [][2]int
	HasAllCards bool
	OK          int
}

func (context *GameContext) getTrade(tradeID int) *GameTrade {
	for _, t := range context.trades {
		if t.ID == tradeID {
			return t
		}
	}
	return nil
}

func (context *GameContext) setupTrade(gives [][2]int, wants [][2]int) error {
	currentPlayer := context.getCurrentPlayer()
	if !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(currentPlayer.ID, gives) {
		return errors.New(ErrInvalidOperation)
	}

	for _, otherPlayer := range context.Players {
		if otherPlayer.ID != currentPlayer.ID {

			hasAllCards := false
			if context.isPlayerHasAllCards(otherPlayer.ID, wants) {
				hasAllCards = true
			}

			var trade = new(GameTrade)
			trade.ID = context.tradeCounter
			trade.From = currentPlayer.ID
			trade.To = otherPlayer.ID
			trade.Gives = gives
			trade.Wants = wants
			trade.HasAllCards = hasAllCards
			trade.OK = 0
			context.trades = append(context.trades, trade)

			//race condition
			context.tradeCounter++
		}
	}

	return nil
}

func (context *GameContext) overrideTrade(playerID, tradeID int, gives [][2]int, wants [][2]int) error {
	var trade = context.getTrade(tradeID)
	if trade == nil || !context.isSafeTrade(gives, wants) || playerID != trade.To || !context.isPlayerHasAllCards(trade.To, gives) {
		return errors.New(ErrInvalidOperation)
	}

	currentPlayer := context.getCurrentPlayer()

	hasAllCards := false
	if context.isPlayerHasAllCards(currentPlayer.ID, wants) {
		hasAllCards = true
	}
	trade.From = trade.To
	trade.To = context.CurrentPlayerID
	trade.Gives = gives
	trade.Wants = wants
	trade.HasAllCards = hasAllCards
	trade.OK = 0

	return nil
}

func (context *GameContext) acceptTrade(playerID, tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 0 || playerID != trade.To {
		return errors.New(ErrInvalidOperation)
	}
	trade.OK = 1
	return nil
}

func (context *GameContext) rejectTrade(playerID, tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 0 || playerID != trade.To {
		return errors.New(ErrInvalidOperation)
	}
	trade.OK = -1
	return nil
}

func (context *GameContext) completeTrade(tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 1 {
		return errors.New(ErrInvalidOperation)
	}

	trade.OK = 2

	player1 := context.Players[trade.From]
	player2 := context.Players[trade.To]

	for _, card := range trade.Gives {
		player1.Cards[card[0]] -= card[1]
		player2.Cards[card[0]] += card[1]
	}

	for _, card := range trade.Wants {
		player1.Cards[card[0]] += card[1]
		player2.Cards[card[0]] -= card[1]
	}

	return nil
}

func (context *GameContext) bankTrade(gives, wants [][2]int) error {
	if context.phase == Phase4 && len(gives) == 1 && len(wants) == 1 && wants[0][1] == 1 && gives[0][1] > 1 {
		currentPlayer := context.getCurrentPlayer()
		if !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(currentPlayer.ID, gives) {
			return errors.New(ErrInvalidOperation)
		}

		banker := context.Bank
		banker.Begin()
		defer banker.Commit()

		wantCardType := wants[0][0]
		wantTradeCount := wants[0][1]
		giveCardType := gives[0][0]
		giveTradeCount := gives[0][1]

		if currentPlayer.ownPort21 || currentPlayer.ownPort31 {
			if currentPlayer.ownPort21 && currentPlayer.ports21[giveCardType] == 1 && giveTradeCount == 2 {
				if err := banker.Set(giveCardType, giveTradeCount); err != nil {
					banker.Rollback()
					return err
				}
				if _, err := banker.Get(wantCardType, wantTradeCount); err != nil {
					banker.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -2
				currentPlayer.Cards[wantCardType]++
			} else if currentPlayer.ownPort31 && giveTradeCount == 3 {
				if err := banker.Set(giveCardType, giveTradeCount); err != nil {
					banker.Rollback()
					return err
				}
				if _, err := banker.Get(wantCardType, wantTradeCount); err != nil {
					banker.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -3
				currentPlayer.Cards[wantCardType]++
			}
		} else {
			if giveTradeCount == 4 {
				if err := banker.Set(giveCardType, giveTradeCount); err != nil {
					banker.Rollback()
					return err
				}
				if _, err := banker.Get(wantCardType, wantTradeCount); err != nil {
					banker.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -4
				currentPlayer.Cards[wantCardType]++
			}
		}

	}

	return errors.New(ErrInvalidOperation)
}