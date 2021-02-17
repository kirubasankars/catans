package main

import "math/rand"

func (context *GameContext) randomPlaceInitialSettlement() error {
	availableIntersections, _ := context.getPossibleSettlementLocations()
	selectedIntersection := availableIntersections[rand.Intn(len(availableIntersections))]
	return context.putSettlement(false, selectedIntersection)
}

func (context *GameContext) randomPlaceInitialRoad() error {
	availableRoads, _ := context.getPossibleRoads()
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	return context.putRoad(false, selectedRoad)
}

func (context *GameContext) randomDiscardCards() {
	for _, player := range context.Players {
		if yes, cardCount := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
			numberOfCardsRemove := cardCount / 2

			for {
				cardId := rand.Intn(5)
				card := player.Cards[cardId]
				if card == 0 {
					if numberOfCardsRemove <= 0 {
						break
					}
					continue
				}

				r := rand.Intn(card) + 1
				if numberOfCardsRemove-r >= 0 {
					player.Cards[cardId] = card - r
					numberOfCardsRemove -= r
				} else {
					player.Cards[cardId] = card - numberOfCardsRemove
					numberOfCardsRemove = 0
				}

				if numberOfCardsRemove <= 0 {
					break
				}
			}
		}
	}
}

func (context *GameContext) randomPlaceRobber() {
	var occupiedIns []int
	for _, player := range context.Players {
		if player.ID == context.CurrentPlayerID {
			continue
		}
		for _, settlement := range player.Settlements {
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
		for _, settlement := range player.Settlements {
			if settlement.Intersection == context.RobberPlacement {
				playerToRob = player.ID
			}
		}
	}

	if playerToRob > -1 {
		context.stealAPlayer(playerToRob)
	}
}
