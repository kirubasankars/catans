package main

import "errors"

func (context *GameContext) playMonopoly(cardType int) error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasMonopoly := false
	for idx, devCard := range currentPlayer.devCards {
		if devCard == DevCardMonopoly {
			hasMonopoly = true
			currentPlayer.devCards = Remove(currentPlayer.devCards, idx)
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
		currentPlayer.cards[cardType] += otherPlayer.cards[cardType]
		otherPlayer.cards[cardType] = 0
	}

	return nil
}