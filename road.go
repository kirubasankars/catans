package main

import (
	"errors"
	"math/rand"
)

func (context *GameContext) getPossibleRoads() ([][2]int, error) {

	if Phase4 == context.phase {
		currentPlayer := context.getCurrentPlayer()

		noOfCoords := len(currentPlayer.Settlements) + len(currentPlayer.Roads)
		var occupiedIns = make([]int, noOfCoords)
		i := 0
		for _, settlement := range currentPlayer.Settlements {
			occupiedIns[i] = settlement.Intersection
			i++
		}
		for _, road := range currentPlayer.Roads {
			if !Contains(occupiedIns, road[0]) {
				occupiedIns[i] = road[0]
			}
			if !Contains(occupiedIns, road[1]) {
				occupiedIns[i] = road[1]
			}
			i++
		}

		i = 0
		var roads = make([][2]int, noOfCoords*3)
		for _, ins := range occupiedIns {
			neighborIntersections := context.board.GetNeighborIntersections2(ins)
			for _, ni := range neighborIntersections {
				roads[i] = ni
				i++
			}
		}

		var allRoads [][2]int
		for _, player := range context.Players {
			allRoads = append(allRoads, player.Roads...)
		}

		var availableRoads [][2]int
		for _, newRoad := range roads {
			found := false
			for _, currentRoad := range allRoads {
				if newRoad[0] == currentRoad[0] && newRoad[1] == currentRoad[1] {
					found = true
					break
				}
			}
			if !found {
				availableRoads = append(availableRoads, newRoad)
			}
		}
		return availableRoads, nil
	}

	if Phase2 == context.phase || Phase3 == context.phase {
		getRoadsForIntersection := func(settlement *Settlement) [][2]int {
			var roads [][2]int
			if settlement != nil {
				neighborIntersections := context.board.GetNeighborIntersections2(settlement.Intersection)
				roads = append(roads, neighborIntersections...)
			}
			return roads
		}

		if Phase2 == context.phase {
			if context.getActionString() != ActionPlaceRoad {
				return nil, errors.New(ErrInvalidOperation)
			}
			currentPlayer := context.getCurrentPlayer()

			var firstSettlement *Settlement
			if len(currentPlayer.Settlements) > 0 {
				firstSettlement = &currentPlayer.Settlements[0]
			}
			return getRoadsForIntersection(firstSettlement), nil
		}

		if Phase3 == context.phase {
			nextAction := context.getAction()
			if nextAction != nil && nextAction.Name != ActionPlaceRoad {
				return nil, errors.New(ErrInvalidOperation)
			}

			currentPlayer := context.getCurrentPlayer()
			var (
				settlementCounter = 1
				secondSettlement  *Settlement
			)

			if len(currentPlayer.Settlements) > 1 {
				settlementCounter++
				secondSettlement = &currentPlayer.Settlements[1]
			}

			return getRoadsForIntersection(secondSettlement), nil
		}
	}

	return nil, errors.New(ErrInvalidOperation)
}

func (context *GameContext) validateRoadPlacement(road [2]int) error {
	availableRoads, _ := context.getPossibleRoads()
	matched := false
	for _, availableRoad := range availableRoads {
		if availableRoad == road {
			matched = true
			break
		}
	}
	if !matched {
		return errors.New(ErrInvalidOperation)
	}
	return nil
}

func (context *GameContext) putRoad(validate bool, road [2]int) error {
	if road[0] > road[1] {
		s := road[1]
		road[1] = road[0]
		road[0] = s
	}
	if validate {
		if err := context.validateRoadPlacement(road); err != nil {
			return err
		}
	}

	if Phase4 == context.phase {
		currentPlayer := context.getCurrentPlayer()
		cards := [][2]int{{0, 1}, {1, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(ErrInvalidOperation)
		}

		banker := context.Bank
		banker.Begin()
		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			if err := banker.Set(card[0], card[1]); err != nil {
				banker.Rollback()
				return err
			}
		}
		banker.Commit()
		currentPlayer.Roads = append(currentPlayer.Roads, road)
		currentPlayer.updateLongestRoad(*context)

		context.EventPutRoad(road)

		return nil
	}

	if Phase2 == context.phase || Phase3 == context.phase {
		if context.getActionString() != ActionPlaceRoad {
			return errors.New(ErrInvalidOperation)
		}
		currentPlayer := context.getCurrentPlayer()
		currentPlayer.Roads = append(currentPlayer.Roads, road)
		currentPlayer.updateLongestRoad(*context)
		context.EventPutRoad(road)
		return context.endAction()
	}
	return nil
}

func (context *GameContext) randomPlaceInitialRoad() error {
	availableRoads, _ := context.getPossibleRoads()
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	return context.putRoad(false, selectedRoad)
}
