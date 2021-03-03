package main

import (
	"errors"
	"math/rand"
)

func (context *GameContext) playKnight() error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer := context.getCurrentPlayer()

	hasPlayKnight := false
	for idx, devCard := range currentPlayer.devCards {
		if devCard == DevCardKnight {
			hasPlayKnight = true
			currentPlayer.devCards = Remove(currentPlayer.devCards, idx)
			break
		}
	}
	if !hasPlayKnight {
		return errors.New(ErrInvalidOperation)
	}

	currentPlayer.knightUsedCount++

	if currentPlayer.knightUsedCount >= 3 {
		for _, otherPlayer := range context.Players {
			if otherPlayer.ID == context.CurrentPlayerID {
				continue
			}

			if currentPlayer.knightUsedCount > otherPlayer.knightUsedCount {
				currentPlayer.hasLargestArmy = true
				currentPlayer.calculateScore()
			}
		}
	}

	context.scheduleAction(ActionPlaceRobber)

	return nil
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
	for _, s := range otherPlayer.settlements {
		if Contains(s.Indices, context.RobberPlacement) {
			hasSettlement = true
		}
	}
	if !hasSettlement {
		return errors.New(ErrInvalidOperation)
	}

	var availableCardTypes []int
	for idx, card := range otherPlayer.cards {
		if card == 0 {
			continue
		}
		availableCardTypes = append(availableCardTypes, idx)
	}

	l := len(availableCardTypes)
	if l > 0 {
		r := rand.Intn(l)
		randCardType := availableCardTypes[r]
		otherPlayer.cards[randCardType]--
		currentPlayer.cards[randCardType]++
	}

	context.scheduleAction(ActionTurn)

	return nil
}

func (context *GameContext) randomPlaceRobber() {
	var occupiedIns []int
	for _, player := range context.Players {
		if player.ID == context.CurrentPlayerID {
			continue
		}
		for _, settlement := range player.settlements {
			occupiedIns = append(occupiedIns, settlement.Intersection)
		}
	}
	ins := rand.Intn(len(occupiedIns))
	tileIndices := context.board.GetTileIndices(ins)
	context.placeRobber(tileIndices[0])
}

func (context *GameContext) randomSelectPlayerToSteal() {
	var playerToRob = -1
	for _, player := range context.Players {
		if player.ID == context.CurrentPlayerID {
			continue
		}
		for _, settlement := range player.settlements {
			for _, tileIndex := range settlement.Indices {
				if tileIndex == context.RobberPlacement {
					playerToRob = player.ID
				}
			}
		}
	}

	if playerToRob == -1 {
		context.scheduleAction(ActionTurn)
		return
	}

	context.stealAPlayer(playerToRob)
}
