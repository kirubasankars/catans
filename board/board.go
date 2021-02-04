package board

import (
	"catans/board/maps"
	"catans/utils"
	"fmt"
	"strings"
)

type Board struct {
	g *grid
}

func (board Board) GetNeighborIntersections1(intersection int) []int {
	thisIntersection := board.g.intersections[intersection]
	neighborIntersections := thisIntersection.neighbors
	var output = make([]int, len(thisIntersection.neighbors))
	for _, ins := range neighborIntersections {
		t := ins.index
		output = append(output, t)
	}
	return output
}

func (board Board) GetNeighborIntersections2(intersection int) [][2]int {
	thisIntersection := board.g.intersections[intersection]
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
	intersections := board.g.intersections
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

func (board Board) GetTileIndices(intersection int) []int {
	coordinator := board.g.intersections[intersection]
	var indices = make([]int, len(coordinator.nodes))
	for idx, n := range coordinator.nodes {
		indices[idx] = n.index
	}
	return indices
}

func (board Board) GetTiles() [][2]int {
	var tiles = make([][2]int, len(board.g.grid))
	var chits = []int{10, 2, 9, 12, 6, 4, 10, 9, 11, 3, 8, 8, 3, 4, 5, 5, 6, 11}
	var tIndex = 0
	for idx, n := range board.g.grid {
		var rt = -1
		var token = -1
		switch n.resource {
		case "t":
			rt = 0
		case "m":
			rt = 1
		case "f":
			rt = 2
		case "p":
			rt = 3
		case "h":
			rt = 4
		default:
			rt = -1
		}
		if rt != -1 {
			token = chits[tIndex]
			tIndex++
		}
		tiles[idx] = [2]int{rt, token}
	}
	return tiles
}

func (board Board) JSON() string {
	tiles := board.GetTiles()
	makeIntersections := func(l []*intersection) string {
		var nodes []string
		for _, h := range l {
			nodes = append(nodes, fmt.Sprintf("{index:%d,x:%f,y:%f, hasPort:%t, portResource:'%s'}", h.index, h.x, h.y, h.hasPort, h.portResource))
		}
		return "[" + strings.Join(nodes, ",") + "]"
	}
	var nodes []string
	for _, h := range board.g.grid {
		if h.resource == "-" || h.resource == "s" {
			continue
		}
		nodes = append(nodes, fmt.Sprintf("{index:%d,x:%f,y:%f,intersections:%s,port:%t,port_direction:%f,resoure:'%s',token:%d}", h.index, h.x, h.y, makeIntersections(h.intersections), h.port, h.direction, h.resource, tiles[h.index][1]))
	}
	return "[" + strings.Join(nodes, ",") + "]"
}

func NewBoard(ID int) Board {
	var g = new(grid)
	if ID == 0 {
		g.Build(maps.DefaultMap{}.GetTileConfig())
	}
	board := new(Board)
	board.g = g
	fmt.Println(board.JSON())
	return *board
}
