package game

import (
	"catans/board"
	"fmt"
	"strings"
)

type Player struct {
	Id          int
	Cards       [5]int
	Roads       [][2]int
	Settlements []Settlement
}

func (player *Player) PutRoad(points [2]int) {
	player.Roads = append(player.Roads, points)
}

func (player *Player) PutSettlement(settlement Settlement) {
	player.Settlements = append(player.Settlements, settlement)
}

func (player Player) Stat() {
	var lines []string
	for idx, count := range player.Cards {
		name := board.ConvertCardTypeToName(idx)
		lines = append(lines, fmt.Sprintf("%s:%d", name, count))
	}
	fmt.Println(player.Id, strings.Join(lines, ", "))
}

func (player Player) String() string {
	return fmt.Sprintf("Player %d", player.Id)
}

func NewPlayer() *Player {
	p := new(Player)
	return p
}
