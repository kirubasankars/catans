package main

import (
	"errors"
	"math/rand"
)

func (context *GameContext) playKnight(tileID, playerID int) error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasPlayKnight := false
	for idx, devCard := range currentPlayer.DevCards {
		if devCard == DevCardKnight {
			hasPlayKnight = true
			currentPlayer.DevCards = Remove(currentPlayer.DevCards, idx)
			break
		}
	}
	if !hasPlayKnight {
		return errors.New(ErrInvalidOperation)
	}

	context.RobberPlacement = tileID
	return context.stealAPlayer(playerID)
}

func (context *GameContext) placeRobber(tileID int) error {
	if context.Action.Name != ActionPlaceRobber {
		return errors.New(ErrInvalidOperation)
	}
	context.RobberPlacement = tileID
	context.scheduleAction(ActionSelectToSteal)
	return nil
}

func (context *GameContext) stealAPlayer(otherPlayerID int) error {
	if context.Action.Name != ActionSelectToSteal {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()
	otherPlayer := context.Players[otherPlayerID]

	// if other player don't have settlement on that tile, throw.
	hasSettlement := false
	for _, s := range otherPlayer.Settlements {
		if Contains(s.Indices, context.RobberPlacement) {
			hasSettlement = true
		}
	}
	if !hasSettlement {
		return errors.New(ErrInvalidOperation)
	}
	// if other player don't have settlement on that tile, throw.

	var availableCards []int
	for idx, card := range otherPlayer.Cards {
		if card == 0 {
			continue
		}
		availableCards = append(availableCards, idx)
	}
	l := len(availableCards)
	if l > 0 {
		r := rand.Intn(l)
		randCardType := availableCards[r]
		otherPlayer.Cards[randCardType] -= 1
		currentPlayer.Cards[randCardType] += 1
	}
	context.scheduleAction(ActionTurn)
	return nil
}
