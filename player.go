package main

import (
	"fmt"
)

type Player struct {
	ID int
	// 0 - lumber
	// 1 - brick
	// 2 - wool
	// 3 - grain
	// 4 - ore

	// 0 - tree
	// 1 - hill
	// 2 - pasture
	// 3 - field
	// 4 - mountain
	Cards          [5]int
	Roads          [][2]int
	Settlements    []Settlement
	DevCards       []int
	hasLargestArmy bool
	hasLongestRoad bool

	ownPort31 bool
	ownPort21 bool
	ports21   [5]int
}

func (player Player) CalculateScore() int {
	score := 0
	for _, settlement := range player.Settlements {
		score += 1
		if settlement.Upgraded {
			score += 1
		}
	}
	for _, devCard := range player.DevCards {
		if devCard == 1 {
			score += 1
		}
	}
	if player.hasLargestArmy {
		score += 2
	}
	if player.hasLongestRoad {
		score += 2
	}
	return score
}

type path struct {
	intersection int
	visited      [][2]int
	length       int
}

func (player Player) uniqueRoadNodes() []int {
	var nodes []int
	for _, road := range player.Roads {
		if !Contains(nodes, road[0]) {
			nodes = append(nodes, road[0])
		}
		if !Contains(nodes, road[1]) {
			nodes = append(nodes, road[1])
		}
	}
	return nodes
}

func (player Player) String() string {
	return fmt.Sprintf("Player %d", player.ID)
}

func (player Player) hasMoreCardsThen(limit int) (bool, int) {
	cardCount := 0
	for _, card := range player.Cards {
		cardCount += card
	}
	return cardCount > limit, cardCount
}

func NewPlayer() *Player {
	p := new(Player)
	return p
}
