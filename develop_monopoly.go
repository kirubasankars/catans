package main

import "errors"

func (context *GameContext) playMonopoly(cardType int) error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasMonopoly := false
	for idx, devCard := range currentPlayer.DevCards {
		if devCard == DevCardMonopoly {
			hasMonopoly = true
			currentPlayer.DevCards = Remove(currentPlayer.DevCards, idx)
			break
		}
	}

	if !hasMonopoly {
		return errors.New(ErrInvalidOperation)
	}

	for _, otherPlayer := range context.Players {
		if otherPlayer.ID == currentPlayer.ID {
			continue
		}
		currentPlayer.Cards[cardType] += otherPlayer.Cards[cardType]
		otherPlayer.Cards[cardType] = 0
	}

	return nil
}