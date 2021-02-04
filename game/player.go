package game

import (
	"catans/utils"
	"fmt"
	"strings"
)

type Player struct {
	ID          int
	Cards       [5]int
	Roads       [][2]int
	Settlements []Settlement
	DevCards    []int

	ownPort31 bool
	ownPort21 bool
	ports21   [5]byte
}

func (player Player) stat() {
	var lines []string
	for idx, count := range player.Cards {
		name := convertCardTypeToTerrain(idx)
		lines = append(lines, fmt.Sprintf("%s:%d", name, count))
	}
	fmt.Println(player.ID, strings.Join(lines, ", "))
}

type path struct {
	intersection int
	visited      [][2]int
	length       int
}

func (player Player) uniqueRoadNodes() []int {
	var nodes []int
	for _, road := range player.Roads {
		if !utils.Contains(nodes, road[0]) {
			nodes = append(nodes, road[0])
		}
		if !utils.Contains(nodes, road[1]) {
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
