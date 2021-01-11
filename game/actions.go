package game

import (
	"catans/dice"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Action struct {
	Name    string
	Timeout time.Time
}

func IsActionTimeout(gc GameContext, action Action) bool {
	timeout := 0
	if gc.GameSetting.Speed == 0 {
		switch action.Name {
		case ActionTurn:
			timeout = 30
		case ActionRollDice:
			timeout = 10
		case ActionPlaceSettlement:
			if gc.Phase == Phase2 || gc.Phase == Phase3 {
				timeout = 12
			}
		case ActionPlaceRoad:
			if gc.Phase == Phase2 || gc.Phase == Phase3 {
				timeout = 15
			}
		case ActionDiscardCards:
			timeout = 15
		case ActionPlaceRobber:
			timeout = 10
		case ActionSelectToRob:
			timeout = 10
		}
	}
	durationSeconds := time.Now().Sub(action.Timeout).Seconds()
	if int(durationSeconds) < timeout {
		return false
	}
	return true
}

func HandleDiscardCards(gc *GameContext) {
	//for _, player := range gc.Players {
	//	if len(player.Cards) > 7 {
	//		fmt.Println(fmt.Sprintf("Player %d looses %d cards.", player.Id, len(player.Cards)/2))
	//	}
	//}
}

func GetAvailableIntersections(gc GameContext) ([]int, error) {

	if Phase4 == gc.Phase {
		//currentPlayer := gc.GetCurrentPlayer()
		//for _, _ := range currentPlayer.Roads {
		//
		//}
	}

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		nextAction := gc.getAction()
		if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceSettlement) {
			return nil, errors.New("invalid action")
		}
		occupied := make([]int, 0)
		for _, player := range gc.Players {
			for _, s := range player.Settlements {
				occupied = append(occupied, s.Intersection)
			}
		}
		return gc.board.GetAvailableIntersections(occupied), nil
	}

	return nil, nil
}

func GetAvailableRoads(gc GameContext) ([][2]int, error) {

	if Phase4 == gc.Phase {
		currentPlayer := gc.getCurrentPlayer()
		var usedIntersections = make(map[int]int)
		for _, settlement := range currentPlayer.Settlements {
			usedIntersections[settlement.Intersection] = 0
		}

		for _, road := range currentPlayer.Roads {
			usedIntersections[road[0]] = 0
			usedIntersections[road[1]] = 0
		}

		var roads [][2]int
		for ins, _ := range usedIntersections {
			neighborIntersections := gc.board.GetNeighborIntersections(ins)
			roads = append(roads, neighborIntersections...)
		}

		var allRoads [][2]int
		for _, player := range gc.Players {
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

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		getRoadsForIntersection := func(settlement *Settlement) [][2]int {
			var roads [][2]int
			if settlement != nil {
				neighborIntersections := gc.board.GetNeighborIntersections(settlement.Intersection)
				roads = append(roads, neighborIntersections...)
			}
			return roads
		}

		if Phase2 == gc.Phase {
			nextAction := gc.getAction()
			if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceRoad) {
				return nil, errors.New("invalid action")
			}
			currentPlayer := gc.getCurrentPlayer()

			var firstSettlement *Settlement
			if len(currentPlayer.Settlements) > 0 {
				firstSettlement = &currentPlayer.Settlements[0]
			}
			return getRoadsForIntersection(firstSettlement), nil
		}

		if Phase3 == gc.Phase {
			nextAction := gc.getAction()
			if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceRoad) {
				return nil, errors.New("invalid action")
			}

			currentPlayer := gc.getCurrentPlayer()
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

	return nil, errors.New("invalid action")
}

func PlaceSettlement(gc *GameContext, validate bool, selectedIntersection int) error {
	if validate {
		availableIntersections, _ := GetAvailableIntersections(*gc)
		matched := false
		for _, availableIntersection := range availableIntersections {
			if availableIntersection == selectedIntersection {
				matched = true
			}
		}
		if !matched {
			return errors.New("invalid action")
		}
	}

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		nextAction := gc.getAction()
		if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceSettlement) {
			return errors.New("invalid action")
		}
		tileIndices := gc.board.GetIndices(selectedIntersection)
		tokens := make([]int, len(tileIndices))
		for i, idx := range tileIndices {
			tokens[i] = gc.tiles[idx].Token
		}
		settlement := Settlement{Indices: tileIndices, Tokens: tokens, Intersection: selectedIntersection}
		gc.putSettlement(settlement)
	}

	return nil
}

func PlaceRoad(gc *GameContext, validate bool, selectedRoad [2]int) error {
	if validate {
		availableRoads, _ := GetAvailableRoads(*gc)
		matched := false
		for _, availableRoad := range availableRoads {
			if availableRoad == selectedRoad {
				matched = true
			}
		}
		if !matched {
			return errors.New("invalid action")
		}
	}

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		nextAction := gc.getAction()
		if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceRoad) {
			return errors.New("invalid action")
		}
		gc.putRoad(selectedRoad)
	}

	return nil
}

func Phase2GetNextAction(gc GameContext) string {
	currentPlayer := gc.getCurrentPlayer()
	settlementCount := len(currentPlayer.Settlements)
	roadCount := len(currentPlayer.Roads)

	if settlementCount == 0 && roadCount == 0 {
		return ActionPlaceSettlement
	}
	if settlementCount == 1 && roadCount == 0 {
		return ActionPlaceRoad
	}
	if settlementCount == 1 && roadCount == 1 {
		return ""
	}
	return ""
}

func Phase3GetNextAction(gc GameContext) string {
	currentPlayer := gc.getCurrentPlayer()
	settlementCount := len(currentPlayer.Settlements)
	roadCount := len(currentPlayer.Roads)

	if settlementCount == 1 && roadCount == 1 {
		return ActionPlaceSettlement
	}
	if settlementCount == 2 && roadCount == 1 {
		return ActionPlaceRoad
	}
	return ""
}

func HandlePlaceInitialSettlement(gc *GameContext) {
	availableIntersections, _ := GetAvailableIntersections(*gc)
	selectedIntersection := availableIntersections[rand.Intn(len(availableIntersections))]
	PlaceSettlement(gc, false, selectedIntersection)
}

func HandlePlaceInitialRoad(gc *GameContext) {
	availableRoads, _ := GetAvailableRoads(*gc)
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	PlaceRoad(gc, false, selectedRoad)
}

func HandleRollDice(gc *GameContext) {
	dice1, dice2 := dice.RandomDice()
	sum := dice1 + dice2

	fmt.Println("Rolled ", sum)

	if sum == 7 {
		gc.scheduleAction(ActionDiscardCards)
		return
	}

	for _, player := range gc.Players {
		for _, settlement := range player.Settlements {
			for idx, token := range settlement.Tokens {
				if token == sum {
					terrain := gc.tiles[settlement.Indices[idx]].Terrain
					count, _ := gc.Bank.Borrow(terrain, 1)
					gc.HandOverCards(player, terrain, count)
				}
			}
		}
	}

	roads, _ := GetAvailableRoads(*gc)
	_ = roads
	for _, player := range gc.Players {
		player.stat()
	}

}
