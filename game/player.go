package game

import (
	"fmt"
	"time"
)

type PlayerAction struct {
	Name    string
	Timeout time.Time
}

type Player struct {
	Id    int
	Cards []string
}

func (player Player) String() string {
	return fmt.Sprintf("Player %d", player.Id)
}

func NewPlayer() *Player {
	p := new(Player)
	return p
}
