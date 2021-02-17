package main

import "errors"

func (context *GameContext) playPlenty(cards [2]int) error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasPlay2Resource := false
	for idx, devCard := range currentPlayer.DevCards {
		if devCard == DevCard2Resource {
			hasPlay2Resource = true
			currentPlayer.DevCards = Remove(currentPlayer.DevCards, idx)
			break
		}
	}

	if !hasPlay2Resource {
		return errors.New(ErrInvalidOperation)
	}

	banker := context.Bank
	banker.Begin()

	for _, cardType := range cards {
		if _, err := banker.Get(cardType, 1); err != nil {
			banker.Rollback()
			return err
		}
		currentPlayer.Cards[cardType]++
	}

	banker.Commit()

	return nil
}