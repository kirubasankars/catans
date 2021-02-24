package main

import "errors"

func (context *GameContext) playPlenty(cards [2]int) error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasPlay2Resource := false
	devCardIndex := 0
	for idx, devCard := range currentPlayer.DevCards {
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

	currentPlayer.Cards[cards[0]]++
	currentPlayer.Cards[cards[1]]++

	if hasPlay2Resource {
		currentPlayer.DevCards = Remove(currentPlayer.DevCards, devCardIndex)
	}

	banker.Commit()

	return nil
}
