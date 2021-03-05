package main

import (
	"errors"
	"math/rand"
)

func (context *GameContext) getPossibleSettlementLocations() ([]int, error) {
	if Phase4 == context.phase {

		var occupiedIns []int
		for _, player := range context.Players {
			for _, settlement := range player.settlements {
				occupiedIns = append(occupiedIns, settlement.Intersection)
			}
		}

		for _, v := range occupiedIns {
			neighborIntersections := context.board.GetNeighborIntersections1(v)
			for _, ni := range neighborIntersections {
				if !Contains(occupiedIns, ni) {
					occupiedIns = append(occupiedIns, ni)
				}
			}
		}

		currentPlayer := context.getCurrentPlayer()
		var roadsIntersections []int
		for _, road := range currentPlayer.roads {
			if !Contains(occupiedIns, road[0]) {
				roadsIntersections = append(roadsIntersections, road[0])
			}
			if !Contains(occupiedIns, road[1]) {
				roadsIntersections = append(roadsIntersections, road[1])
			}
		}

		var availableIntersections []int
		for _, ins := range roadsIntersections {
			if !Contains(occupiedIns, ins) {
				availableIntersections = append(availableIntersections, ins)
			}
		}
		return availableIntersections, nil
	}

	if Phase2 == context.phase || Phase3 == context.phase {
		if context.getActionString() != ActionPlaceSettlement {
			return nil, errors.New(ErrInvalidOperation)
		}
		occupied := make([]int, 0)
		for _, player := range context.Players {
			for _, s := range player.settlements {
				occupied = append(occupied, s.Intersection)
			}
		}
		return context.board.GetAvailableIntersections(occupied), nil
	}

	return nil, errors.New(ErrInvalidOperation)
}

func (context *GameContext) putSettlement(validate bool, intersection int) error {
	currentPlayer := context.getCurrentPlayer()
	if context.getActionString() != ActionPlaceSettlement || currentPlayer.allowedSettlementCount < 1 {
		return errors.New(ErrInvalidOperation)
	}

	if validate {
		availableIntersections, _ := context.getPossibleSettlementLocations()
		if !Contains(availableIntersections, intersection) {
			return errors.New(ErrInvalidOperation)
		}
	}

	if Phase2 == context.phase || Phase3 == context.phase {
		var settlement Settlement

		tileIndices := context.board.GetTileIndices(intersection)

		tokens := make([]int, len(tileIndices))
		for i, idx := range tileIndices {
			tokens[i] = context.Tiles[idx][1]
		}

		settlement = Settlement{Indices: tileIndices, Tokens: tokens, Intersection: intersection}
		currentPlayer.settlements = append(currentPlayer.settlements, settlement)
		currentPlayer.allowedSettlementCount--
	}

	if Phase4 == context.phase {
		cards := [][2]int{{CardLumber, 1}, {CardBrick, 1}, {CardWool, 1}, {CardGrain, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(ErrInvalidOperation)
		}

		banker := context.Bank
		banker.Begin()

		for _, card := range cards {
			currentPlayer.cards[card[0]] -= card[1]
			if err := banker.Set(card[0], card[1]); err != nil {
				banker.Rollback()
				return err
			}
		}
		for _, card := range cards {
			currentPlayer.cards[card[0]] -= card[1]
		}

		var settlement Settlement
		tileIndices := context.board.GetTileIndices(intersection)
		tokens := make([]int, len(tileIndices))
		for i, idx := range tileIndices {
			tokens[i] = context.Tiles[idx][1]
		}

		settlement = Settlement{Indices: tileIndices, Tokens: tokens, Intersection: intersection}
		currentPlayer.settlements = append(currentPlayer.settlements, settlement)
		currentPlayer.allowedSettlementCount--

		banker.Commit()
	}

	context.EventPutSettlement(intersection)
	currentPlayer.calculateScore()

	if Phase2 == context.phase || Phase3 == context.phase {
		return context.endAction()
	}

	return nil
}

func (context *GameContext) upgradeSettlement(intersection int) error {
	currentPlayer := context.getCurrentPlayer()
	if Phase4 == context.phase {

		if currentPlayer.allowedSettlementUpgradeCount < 1 {
			return errors.New(ErrInvalidOperation)
		}

		var settlement *Settlement
		for _, s := range currentPlayer.settlements {
			if s.Intersection == intersection {
				settlement = &s
				break
			}
		}
		if settlement == nil {
			return errors.New(ErrInvalidOperation)
		}

		cards := [][2]int{{3, 2}, {4, 3}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(ErrInvalidOperation)
		}

		bank := context.Bank
		bank.Begin()
		for _, card := range cards {
			currentPlayer.cards[card[0]] -= card[1]
			err := bank.Set(card[0], card[1])
			if err != nil {
				bank.Rollback()
				return err
			}
		}

		settlement.Upgraded = true
		currentPlayer.allowedSettlementUpgradeCount--
		currentPlayer.calculateScore()

		for _, card := range cards {
			currentPlayer.cards[card[0]] -= card[1]
		}
		bank.Commit()

		return nil
	}
	return errors.New(ErrInvalidOperation)
}

func (context *GameContext) randomPlaceInitialSettlement() error {
	availableIntersections, _ := context.getPossibleSettlementLocations()
	selectedIntersection := availableIntersections[rand.Intn(len(availableIntersections))]
	return context.putSettlement(false, selectedIntersection)
}
