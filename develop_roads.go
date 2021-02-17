package main

import "errors"

func (context *GameContext) playRoads(roads [][2]int) error {
	if context.phase != Phase4 {
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

	for _, road := range roads {
		if road[0] > road[1] {
			s := road[1]
			road[1] = road[0]
			road[0] = s
		}
		if err := context.validateRoadPlacement(road); err != nil {
			return err
		}
	}
	for _, road := range roads {
		currentPlayer.Roads = append(currentPlayer.Roads, road)
	}
	return nil
}

