package board

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Board struct {
	_map 	Map
}

func (b Board) Test() {
	for _, c := range b._map.coordinators {
		fmt.Println(c, len(c.neighbors))
	}
}

func (b Board) GetNeighborIntersections(intersection string) []string {
	neighborIntersections := b._map.coordinators[intersection].neighbors
	var output []string
	for _, ins := range neighborIntersections {
		output = append(output, ins.String())
	}
	return output
}

func (b Board) GetAvailableIntersections(occupied []string) []string {
	intersections := b._map.coordinators
	l := len(occupied)
	for i := 0; i < l; i++ {
		occupiedNeighbors := intersections[occupied[i]].neighbors
		for _, nins := range occupiedNeighbors {
			occupied = append(occupied, nins.String())
		}
	}
	keys := make([]string, 0, len(intersections))
	for k := range intersections {
		if !contains(occupied, k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (b Board) GetTileIndex(intersection string) []int {
	coord := b._map.coordinators[intersection]
	var indices []int
	for _, n := range coord.nodes {
		indices = append(indices, n.index)
	}
	return indices
}

func (b Board) GenerateTiles(tileSettings string) []Tile {
	r := regexp.MustCompile(`(?P<Token>\d+)(?P<Terrain>\w+)?`)
	segs := strings.Split(tileSettings, ",")
	tiles := make([]Tile, len(segs))
	for idx, seg := range segs {
		rs := r.FindAllStringSubmatch(seg, -1)
		if len(rs) > 0 {
			tiles[idx].Token, _ = strconv.Atoi(rs[0][1])
			if len(rs[0]) > 1 {
				tiles[idx].Terrain = rs[0][2]
			}
		}
	}
	return tiles
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NewBoard() Board {
	return Board{ _map: GetDefaultMap() }
}