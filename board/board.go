package board

import (
	"catans/maps"
	"catans/utils"
)

type Board struct {
	_map 	Map
}

func (board Board) GetNeighborIntersections(intersection int) [][2]int {
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

func (board Board) GetIndices(intersection int) []int {
	coordinator := board._map.coordinators[intersection]
	var indices = make([]int, len(coordinator.nodes))
	for _, n := range coordinator.nodes {
		indices = append(indices, n.index)
	}
	return indices
}

func NewBoard(name string) Board {
	_map := newMap()
	if name == "default" {
		_map.build(maps.DefaultMap{})
	}
	if name == "diamond" {
		_map.build(maps.Diamond{})
	}
	board := new (Board)
	board._map = _map 
	return *board
}