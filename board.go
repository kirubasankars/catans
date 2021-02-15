package main

type Board struct {
	grid *Grid
}

func (board Board) GetNeighborIntersections1(intersection int) []int {
	thisIntersection := board.grid.intersections[intersection]
	neighborIntersections := thisIntersection.neighbors
	var output = make([]int, len(thisIntersection.neighbors))
	for _, ins := range neighborIntersections {
		t := ins.index
		output = append(output, t)
	}
	return output
}

func (board Board) GetNeighborIntersections2(intersection int) [][2]int {
	thisIntersection := board.grid.intersections[intersection]
	neighborIntersections := thisIntersection.neighbors
	var output = make([][2]int, len(thisIntersection.neighbors))
	for _, ins := range neighborIntersections {
		if ins.index < intersection {
			output = append(output, [2]int{ins.index, intersection})
		} else {
			output = append(output, [2]int{intersection, ins.index})
		}
	}
	return output
}

func (board Board) GetAvailableIntersections(occupied []int) []int {
	intersections := board.grid.intersections
	l := len(occupied)
	for i := 0; i < l; i++ {
		occupiedNeighbors := intersections[occupied[i]].neighbors
		for _, nins := range occupiedNeighbors {
			occupied = append(occupied, nins.index)
		}
	}
	keys := make([]int, 0, len(intersections))
	for k := range intersections {
		if !Contains(occupied, k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (board Board) GetTileIndices(intersection int) []int {
	coordinator := board.grid.intersections[intersection]
	var indices = make([]int, len(coordinator.nodes))
	for idx, n := range coordinator.nodes {
		indices[idx] = n.index
	}
	return indices
}

func (board Board) GetTiles() [][2]int {
	var tiles = make([][2]int, len(board.grid.nodes))
	for idx, n := range board.grid.nodes {
		var rt = -1
		switch n.terrain {
		case "t":
			rt = 0
		case "h":
			rt = 1
		case "p":
			rt = 2
		case "f":
			rt = 3
		case "m":
			rt = 4
		case "l":
			rt = 0
		case "b":
			rt = 1
		case "w":
			rt = 2
		case "g":
			rt = 3
		case "o":
			rt = 4
		case "?":
			rt = 5
		case "-":
			rt = -1
		case "s":
			rt = -2
		case "d":
			rt = -3
		}
		tiles[idx] = [2]int{rt, n.token}
	}
	return tiles
}

func NewBoard(ID int) Board {
	var grid = new(Grid)
	if ID == 0 {
		m := DefaultMap{}
		grid.Build(m)
	}
	if ID == 1 {
		m := SmallMap{}
		grid.Build(m)
	}
	board := new(Board)
	board.grid = grid
	return *board
}

func convertCardTypeToTerrain(cardType int) string {
	switch cardType {
	case 0:
		return "t"
	case 1:
		return "h"
	case 2:
		return "p"
	case 3:
		return "f"
	case 4:
		return "m"
	case 5:
		return "?"
	case -1:
		return "-"
	case -2:
		return "s"
	case -3:
		return "d"
	}
	return ""
}
