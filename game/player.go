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

	has31   bool
	has21   bool
	cards21 [5]byte
}

func (player *Player) putRoad(points [2]int) error {
	player.Roads = append(player.Roads, points)
	return nil
}

func (player *Player) putSettlement(settlement Settlement) error {
	player.Settlements = append(player.Settlements, settlement)
	return nil
}

func (player Player) stat() {
	var lines []string
	for idx, count := range player.Cards {
		name := convertCardTypeToName(idx)
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

func NewPlayer() *Player {
	p := new(Player)
	return p
}
