package main

import (
	"errors"
	"math/rand"
)

func (context *GameContext) playRoads() error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}

	availableRoads, _ := context.getPossibleRoads()
	if len(availableRoads) == 0 {
		return errors.New(ErrInvalidOperation)
	}

	currentPlayer := context.getCurrentPlayer()
	hasPlay2Road := false
	for idx, devCard := range currentPlayer.DevCards {
		if devCard == DevCard2Road {
			hasPlay2Road = true
			currentPlayer.DevCards = Remove(currentPlayer.DevCards, idx)
			break
		}
	}

	if !hasPlay2Road {
		return errors.New(ErrInvalidOperation)
	}

	context.scheduleAction(ActionDevPlaceRoad1)

	return nil
}

func (context *GameContext) playDev2PlaceRoad(road [2]int) error {
	if !(context.Action.Name == ActionDevPlaceRoad1 || context.Action.Name == ActionDevPlaceRoad2) {
		return errors.New(ErrInvalidOperation)
	}

	if road[0] > road[1] {
		s := road[1]
		road[1] = road[0]
		road[0] = s
	}
	if err := context.validateRoadPlacement(road); err != nil {
		return err
	}

	currentPlayer := context.getCurrentPlayer()
	currentPlayer.Roads = append(currentPlayer.Roads, road)

	updateLongestRoad := func() {
		if currentPlayer.RoadsCount > 4 {
			for _, otherPlayer := range context.Players {
				if otherPlayer.ID == context.CurrentPlayerID {
					continue
				}

				if currentPlayer.RoadsCount > otherPlayer.RoadsCount {
					currentPlayer.hasLongestRoad = true
					currentPlayer.CalculateScore()
				}
			}
		}
	}

	currentPlayer.RoadsCount++
	updateLongestRoad()

	if context.Action.Name == ActionDevPlaceRoad1 {
		availableRoads, _ := context.getPossibleRoads()
		if len(availableRoads) == 0 {
			currentPlayer.Roads = currentPlayer.Roads[:len(currentPlayer.Roads)-1]
			currentPlayer.DevCards = append(currentPlayer.DevCards, DevCard2Road)
			currentPlayer.RoadsCount--
			updateLongestRoad()
			context.scheduleAction(ActionTurn)
			return errors.New(ErrInvalidOperation)
		}
		context.scheduleAction(ActionDevPlaceRoad2)
	}

	if context.Action.Name == ActionDevPlaceRoad2 {
		context.scheduleAction(ActionTurn)
	}

	return nil
}

func (context *GameContext) randomPlaceDev2PlaceRoad() {
	availableRoads, _ := context.getPossibleRoads()
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	context.playDev2PlaceRoad(selectedRoad)
}