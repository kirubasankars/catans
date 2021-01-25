package board

import (
	"catans/maps"
	"catans/utils"
	"fmt"
	"strings"
)

type Board struct {
	_map Map
}

func (board Board) GetNeighborIntersections1(intersection int) []int {
	thisIntersection := board._map.coordinators[intersection]
	neighborIntersections := thisIntersection.neighbors
	var output []int
	for _, ins := range neighborIntersections {
		t := ins.index
		output = append(output, t)
	}
	return output
}

func (board Board) GetNeighborIntersections2(intersection int) [][2]int {
	thisIntersection := board._map.coordinators[intersection]
	neighborIntersections := thisIntersection.neighbors
	var output [][2]int
	for _, ins := range neighborIntersections {
		t := ins.index
		if thisIntersection.nodes[0].index == ins.nodes[0].index {
			if thisIntersection.sides[0] < ins.sides[0] {
				output = append(output, [2]int{intersection, t})
			} else {
				output = append(output, [2]int{t, intersection})
			}
		} else if thisIntersection.nodes[0].index < ins.nodes[0].index {
			output = append(output, [2]int{intersection, t})
		} else {
			output = append(output, [2]int{t, intersection})
		}
	}
	return output
}

func (board Board) GetAvailableIntersections(occupied []int) []int {
	intersections := board._map.coordinators
	l := len(occupied)
	for i := 0; i < l; i++ {
		occupiedNeighbors := intersections[occupied[i]].neighbors
		for _, nins := range occupiedNeighbors {
			occupied = append(occupied, nins.index)
		}
	}
	keys := make([]int, 0, len(intersections))
	for k := range intersections {
		if !utils.Contains(occupied, k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (board Board) GetTiles(intersection int) []int {
	coordinator := board._map.coordinators[intersection]
	var indices = make([]int, len(coordinator.nodes))
	for idx, n := range coordinator.nodes {
		indices[idx] = n.index
	}
	return indices
}

func CatNodes(ID int) string {
	b := NewBoard(ID)
	var output []string
	for i := 0; i < len(b._map.nodes); i++ {
		output = append(output, b._map.nodes[i].String())
	}
	return strings.Join(output, "\n")
}

func CatIntersections(ID int) string {
	b := NewBoard(ID)
	var output []string
	for i := 0; i < len(b._map.coordinators); i++ {
		nc := b._map.coordinators[i]
		output = append(output, fmt.Sprint(nc.index, nc, nc.neighbors))
	}
	return strings.Join(output, "\n")
}

func NewBoard(ID int) Board {
	_map := newMap()
	if ID == 0 {
		_map.build(maps.DefaultMap{})
	}
	if ID == 1 {
		_map.build(maps.Diamond{})
	}
	board := new(Board)
	board._map = _map
	return *board
}
