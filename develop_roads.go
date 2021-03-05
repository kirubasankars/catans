package main

import (
	"errors"
	"math/rand"
)

func (context *GameContext) playRoads() error {
	if context.phase != Phase4 {
		return errors.New(ErrInvalidOperation)
	}

	currentPlayer := context.getCurrentPlayer()
	availableRoads, _ := context.getPossibleRoads()
	if len(availableRoads) == 0 || currentPlayer.allowedRoadsCount < 1 {
		return errors.New(ErrInvalidOperation)
	}


	hasPlay2Road := false
	for idx, devCard := range currentPlayer.devCards {
		if devCard == DevCard2Road {
			hasPlay2Road = true
			currentPlayer.devCards = Remove(currentPlayer.devCards, idx)
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
	if currentPlayer.allowedRoadsCount < 1 {
		return errors.New(ErrInvalidOperation)
	}
	currentPlayer.roads = append(currentPlayer.roads, road)
	currentPlayer.allowedRoadsCount--
	currentPlayer.updateLongestRoad(*context)

	if context.Action.Name == ActionDevPlaceRoad1 {
		availableRoads, _ := context.getPossibleRoads()
		if len(availableRoads) == 0 {

			currentPlayer.roads = currentPlayer.roads[:len(currentPlayer.roads)-1]
			currentPlayer.allowedRoadsCount++
			currentPlayer.devCards = append(currentPlayer.devCards, DevCard2Road)
			currentPlayer.updateLongestRoad(*context)

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

func (context *GameContext) randomPlaceDev2PlaceRoad() error {
	availableRoads, _ := context.getPossibleRoads()
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	context.playDev2PlaceRoad(selectedRoad)
	return nil
}