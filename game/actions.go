package game

import (
	"catans/dice"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func IsActionTimeout(gc GameContext, action PlayerAction) bool {
	timeout := 0
	if gc.setting.Speed == 0 {
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

func TerrainCard(terrain string) string {
	switch terrain {
	case "MO":
		return "S"
	case "FO":
		return "L"
	case "HI":
		return "B"
	case "FI":
		return "G"
	case "PA":
		return "W"
	}
	return ""
}

func HandleDiscardCards(gc *GameContext) {
	for _, player := range gc.Players {
		if len(player.Cards) > 7 {
			fmt.Println(fmt.Sprintf("Player %d looses %d cards.", player.Id, len(player.Cards)/2))
		}
	}
}

func GetAvailableIntersections(gc GameContext) ([]string, error) {

	if Phase4 == gc.Phase {

	}

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		nextAction := gc.GetAction()
		if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceSettlement) {
			return nil, errors.New("invalid action")
		}
		occupied := make([]string, 0, len(gc.Settlements))
		for _, s := range gc.Settlements {
			occupied = append(occupied, s.Intersection)
		}
		return gc.board.GetAvailableIntersections(occupied), nil
	}

	return nil, nil
}

func GetAvailableRoads(gc GameContext) ([][2]string, error) {

	if Phase4 == gc.Phase {
		currentPlayer := gc.GetCurrentPlayer()
		var roadIntersections []string
		for _, road := range gc.Roads {
			if road.player == currentPlayer {
				roadIntersections = append(roadIntersections, road.Points[0], road.Points[1])
			}
		}
		var settlementIntersections []string
		for _, road := range gc.Settlements {
			if road.player == currentPlayer {
				settlementIntersections = append(settlementIntersections, road.Intersection)
			}
		}

		//remove settlement intersections from road
		for i := 0; i < len(roadIntersections); i++ {
			for _, intersection := range settlementIntersections {
				if intersection == roadIntersections[i] {

				}
			}
		}

	}

	if Phase2 == gc.Phase || Phase3 == gc.Phase {
		getRoadsForIntersection := func(settlement *Settlement) [][2]string {
			var roads [][2]string
			if settlement != nil {
				neighborIntersections := gc.board.GetNeighborIntersections(settlement.Intersection)
				for _, ins := range neighborIntersections {
					roads = append(roads, [2]string{settlement.Intersection, ins})
				}
			}
			return roads
		}

		if Phase2 == gc.Phase {
			nextAction := gc.GetAction()
			if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceRoad) {
				return nil, errors.New("invalid action")
			}
			currentPlayer := gc.GetCurrentPlayer()
			var firstSettlement *Settlement
			for _, settlement := range gc.Settlements {
				if settlement.player == currentPlayer {
					firstSettlement = &settlement
				}
			}
			return getRoadsForIntersection(firstSettlement), nil
		}

		if Phase3 == gc.Phase {
			nextAction := gc.GetAction()
			if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceRoad) {
				return nil, errors.New("invalid action")
			}
			currentPlayer := gc.GetCurrentPlayer()
			var (
				settlementCounter = 1
				secondSettlement  *Settlement
			)

			for _, settlement := range gc.Settlements {
				if settlement.player == currentPlayer {
					if settlementCounter == 2 {
						secondSettlement = &settlement
					}
					settlementCounter++
				}
			}

			return getRoadsForIntersection(secondSettlement), nil
		}
	}

	return nil, errors.New("invalid action")
}

func PlaceSettlement(gc *GameContext, validate bool, selectedIntersection string) error {
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
		nextAction := gc.GetAction()
		if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceSettlement) {
			return errors.New("invalid action")
		}
		tileIndices := gc.board.GetTileIndex(selectedIntersection)
		tokens := make([]int, len(tileIndices))
		for i, idx := range tileIndices {
			tokens[i] = gc.tiles[idx].Token
		}
		currentPlayer := gc.GetCurrentPlayer()
		settlement := Settlement{ player: currentPlayer, Indices: tileIndices, Tokens: tokens, Intersection: selectedIntersection}
		gc.PutSettlement(settlement)
	}

	return nil
}

func PlaceRoad(gc *GameContext, validate bool, selectedRoad [2]string) error {
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
		nextAction := gc.GetAction()
		if nextAction == nil || (nextAction != nil && nextAction.Name != ActionPlaceRoad) {
			return errors.New("invalid action")
		}
		gc.PutRoad(Road{player: gc.GetCurrentPlayer(), Points: selectedRoad})
	}

	return nil
}

func Phase2GetNextAction(gc GameContext) string {
	currentPlayer := gc.GetCurrentPlayer()
	settlementCount := 0
	for _, settlement := range gc.Settlements {
		if settlement.player == currentPlayer {
			settlementCount ++
		}
	}

	roadCount := 0
	for _, road := range gc.Roads {
		if road.player == currentPlayer {
			roadCount ++
		}
	}

	if settlementCount == 0 && roadCount == 0  {
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
	currentPlayer := gc.GetCurrentPlayer()
	settlementCount := 0
	for _, settlement := range gc.Settlements {
		if settlement.player == currentPlayer {
			settlementCount ++
		}
	}

	roadCount := 0
	for _, road := range gc.Roads {
		if road.player == currentPlayer {
			roadCount ++
		}
	}

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
		gc.ScheduleAction(ActionDiscardCards)
	}

	for _, settlement := range gc.Settlements {
		var cards []string
		for idx, token := range settlement.Tokens {
			if token == sum {
				terrain := gc.tiles[settlement.Indices[idx]].Terrain
				cards = append(cards, TerrainCard(terrain))
			}
		}
		if len(cards)> 0 {
			gc.HandOverCards(settlement.player, cards)
			PrintStat(gc)
		}
	}
}

func PrintStat(gc *GameContext) {
	for _, player := range gc.Players {
		for _, settlement := range gc.Settlements {
			if settlement.player.Id == player.Id {
				//fmt.Println(settlement.Intersection)
			}
		}
		var text []string
		var cardstat = make(map[string]int)
		for _, card := range player.Cards {
			cardstat[card]++
		}
		for k, v := range cardstat {
			text = append(text, fmt.Sprintf("%d%s", v, k))
		}
		if len(cardstat) > 0 {
			fmt.Println(player.Id, text)
		}
	}
}
