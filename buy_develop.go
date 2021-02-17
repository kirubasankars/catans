package main

import "errors"

func (context *GameContext) buyDevelopmentCard() error {
	currentPlayer := context.getCurrentPlayer()
	if Phase4 == context.phase {
		cards := [][2]int{{2, 1}, {3, 1}, {4, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(ErrInvalidOperation)
		}

		bank := context.Bank
		bank.Begin()

		card, err := context.Bank.BuyDevCard()
		if err != nil {
			bank.Rollback()
			return err
		}

		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			err := context.Bank.Set(card[0], card[1])
			if err != nil {
				bank.Rollback()
				return err
			}
		}

		currentPlayer.DevCards = append(currentPlayer.DevCards, card)

		context.Bank.Commit()
	}
	return errors.New(ErrInvalidOperation)
}
