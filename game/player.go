package game

import (
	"catans/utils"
	"container/list"
	"fmt"
	"strings"
)

type Player struct {
	ID          int
	Cards       [5]int
	Roads       [][2]int
	Settlements []Settlement
}

func (player *Player) putRoad(points [2]int) {
	player.Roads = append(player.Roads, points)
}

func (player *Player) putSettlement(settlement Settlement) {
	player.Settlements = append(player.Settlements, settlement)
}

func (player Player) stat() {
	var lines []string
	for idx, count := range player.Cards {
		name := ConvertCardTypeToName(idx)
		lines = append(lines, fmt.Sprintf("%s:%d", name, count))
	}
	fmt.Println(player.ID, strings.Join(lines, ", "))
}

type path struct {
	intersection int
	visited [][2]int
	length  int
}

func (player Player) calculateLongestRoad(otherPlayersSettlements []int) int {
	roadNodes := player.uniqueRoadNodes()
	pending := list.New()
	longest := 0

	for _, node := range roadNodes {

		pending.PushBack(path{ intersection: node, length: 0, visited: [][2]int{}})

		fmt.Println("FROM", node)

		for pending.Len() > 0 {

			var pathEnd = true

			el := pending.Front()
			pending.Remove(el)

			item := el.Value.(path)

			for _, road := range player.Roads {

				if road[0] == item.intersection || road[1] == item.intersection {
						////broken road check
						//if otherPlayersSettlements != nil && utils.Contains(otherPlayersSettlements, r1) {
						//	pathEnd = true
						//}

						p := -1
						if road[0] == item.intersection {
							p = road[1]
						} else {
							p = road[0]
						}


						visited := false
						for _, v := range item.visited {
							if v[0] == road[0] && v[1] == road[1] {
								visited = true
							}
						}
						if !visited {
							pathEnd = false
							item.visited = append(item.visited, road)
							pending.PushBack(path{intersection: p, length: item.length + 1, visited: item.visited})
						}

				}
			}

			if pathEnd {
				fmt.Println("Path End", item)
				if longest < item.length {
					longest = item.length
				}
			}
		}
	}

	return longest
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
