package main

import "errors"

func (context *GameContext) buyDevelopmentCard() error {
	currentPlayer := context.getCurrentPlayer()
	if Phase4 == context.phase {
		cards := [][2]int{{CardWool, 1}, {CardGrain, 1}, {CardOre, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(ErrInvalidOperation)
		}

		bank := context.Bank
		bank.Begin()

		card, err := bank.BuyDevCard()
		if err != nil {
			bank.Rollback()
			return err
		}

		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			err := bank.Set(card[0], card[1])
			if err != nil {
				bank.Rollback()
				return err
			}
		}

		currentPlayer.DevCards = append(currentPlayer.DevCards, card)

		bank.Commit()
	}
	return errors.New(ErrInvalidOperation)
}
