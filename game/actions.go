package game

import (
	"catans/utils"
	"errors"
	"fmt"
	"math/rand"
)

func GetPossibleSettlementLocations(gc GameContext) ([]int, error) {
	if Phase4 == gc.Phase {

		var occupiedIns []int
		for _, player := range gc.Players {
			for _, settlement := range player.Settlements {
				occupiedIns = append(occupiedIns, settlement.Intersection)
			}
		}

		for _, v := range occupiedIns {
			neighborIntersections := gc.board.GetNeighborIntersections1(v)
			for _, ni := range neighborIntersections {
				if !utils.Contains(occupiedIns, ni) {
					occupiedIns = append(occupiedIns, ni)
				}
			}
		}

		currentPlayer := gc.getCurrentPlayer()
		var roadsIntersections []int
		for _, road := range currentPlayer.Roads {
			if !utils.Contains(occupiedIns, road[0]) {
				roadsIntersections = append(roadsIntersections, road[0])
			}
			if !utils.Contains(occupiedIns, road[1]) {
				roadsIntersections = append(roadsIntersections, road[1])
			}
		}

		var availableIntersections []int
		for _, ins := range roadsIntersections {
			if !utils.Contains(occupiedIns, ins) {
				availableIntersections = append(availableIntersections, ins)
			}
		}
		return availableIntersections, nil
	}

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		nextAction := gc.getAction()
		if nextAction != nil && nextAction.Name != ActionPlaceSettlement {
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

func GetPossibleRoads(gc GameContext) ([][2]int, error) {

	if Phase4 == gc.Phase {
		currentPlayer := gc.getCurrentPlayer()

		noOfCoords := len(currentPlayer.Settlements) + len(currentPlayer.Roads)
		var occupiedIns = make([]int, noOfCoords)
		i := 0
		for _, settlement := range currentPlayer.Settlements {
			occupiedIns[i] = settlement.Intersection
			i++
		}
		for _, road := range currentPlayer.Roads {
			if !utils.Contains(occupiedIns, road[0]) {
				occupiedIns[i] = road[0]
			}
			if !utils.Contains(occupiedIns, road[1]) {
				occupiedIns[i] = road[1]
			}
			i++
		}

		i = 0
		var roads = make([][2]int, noOfCoords*3)
		for _, ins := range occupiedIns {
			neighborIntersections := gc.board.GetNeighborIntersections2(ins)
			for _, ni := range neighborIntersections {
				roads[i] = ni
				i++
			}
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
				neighborIntersections := gc.board.GetNeighborIntersections2(settlement.Intersection)
				roads = append(roads, neighborIntersections...)
			}
			return roads
		}

		if Phase2 == gc.Phase {
			nextAction := gc.getAction()
			if nextAction != nil && nextAction.Name != ActionPlaceRoad {
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
			if nextAction != nil && nextAction.Name != ActionPlaceRoad {
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
		availableIntersections, _ := GetPossibleSettlementLocations(*gc)
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
		if nextAction != nil && nextAction.Name != ActionPlaceSettlement {
			return errors.New("invalid action")
		}
		tileIndices := gc.board.GetTiles(selectedIntersection)
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
		availableRoads, _ := GetPossibleRoads(*gc)
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
		if nextAction != nil && nextAction.Name != ActionPlaceRoad {
			return errors.New("invalid action")
		}
		gc.putRoad(selectedRoad)
	}

	return nil
}

func RandomPlaceInitialSettlement(context *GameContext) {
	availableIntersections, _ := GetPossibleSettlementLocations(*context)
	selectedIntersection := availableIntersections[rand.Intn(len(availableIntersections))]
	PlaceSettlement(context, false, selectedIntersection)
}

func RandomPlaceInitialRoad(context *GameContext) {
	availableRoads, _ := GetPossibleRoads(*context)
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	PlaceRoad(context, false, selectedRoad)
}

func HandleDice(context *GameContext, dice int) error {
	if dice == 7 {
		return errors.New("invalid operation")
	}

	bank := context.Bank
	players := context.Players

	bank.Begin()
	for _, player := range players {
		for _, settlement := range player.Settlements {
			for idx, token := range settlement.Tokens {
				if token == dice {
					terrain := context.tiles[settlement.Indices[idx]].Terrain
					count, err := context.Bank.Give(terrain, 1)
					if err != nil {
						bank.Rollback()
						return err
					}
					context.handOverCards(player, terrain, count)
				}
			}
		}
	}

	bank.Commit()

	return nil
}

func RandomDiscardCards(gc *GameContext) {
	for _, player := range gc.Players {
		if len(player.Cards) > 7 {
			fmt.Println(fmt.Sprintf("Player %d looses %d cards.", player.ID, len(player.Cards)/2))
		}
	}
}