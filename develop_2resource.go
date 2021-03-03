package main

import "errors"

func (context *GameContext) play2Resource(cards [2]int) error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasPlay2Resource := false
	devCardIndex := 0
	for idx, devCard := range currentPlayer.devCards {
		if devCard == DevCard2Resource {
			hasPlay2Resource = true
			devCardIndex = idx
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
	}

	currentPlayer.cards[cards[0]]++
	currentPlayer.cards[cards[1]]++

	if hasPlay2Resource {
		currentPlayer.devCards = Remove(currentPlayer.devCards, devCardIndex)
	}

	banker.Commit()

	return nil
}
