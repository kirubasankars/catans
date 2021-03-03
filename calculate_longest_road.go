package main

import "container/list"

func (context GameContext) calculateLongestRoad(player Player, otherPlayersSettlements []int) int {

	uniqueRoads := func() []int {
		var nodes []int
		for _, road := range player.roads {
			if !Contains(nodes, road[0]) {
				nodes = append(nodes, road[0])
			}
			if !Contains(nodes, road[1]) {
				nodes = append(nodes, road[1])
			}
		}
		return nodes
	}

	roadNodes := uniqueRoads()
	pending := list.New()
	longest := 0

	for _, node := range roadNodes {

		pending.PushBack(path{intersection: node, length: 0, visited: [][2]int{}})

		//fmt.Println("FROM", node)

		for pending.Len() > 0 {

			var pathEnd = true

			el := pending.Front()
			pending.Remove(el)

			item := el.Value.(path)

			for _, road := range player.roads {

				if road[0] == item.intersection || road[1] == item.intersection {
					////broken road check
					//if otherPlayersSettlements != nil && Contains(otherPlayersSettlements, r1) {
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
				//fmt.Println("Path End", item)
				if longest < item.length {
					longest = item.length
				}
			}
		}
	}

	return longest
}

type path struct {
	intersection int
	visited      [][2]int
	length       int
}
