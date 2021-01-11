package board

import (
	"regexp"
	"strconv"
	"strings"
)

func ConvertTerrainToCardType(terrain string) int {
	switch terrain {
	case "FO":
		return  0
	case "HI":
		return  1
	case "PA":
		return  2
	case "FI":
		return  3
	case "MO":
		return 4
	}
	return -1
}

func ConvertCardTypeToName(cardType int) string {
	switch cardType {
	case 0:
		return "Log"
	case 1:
		return "Brick"
	case 2:
		return "Wool"
	case 3:
		return "Grain"
	case 4:
		return "Ore"
	}
	return ""
}

type Board struct {
	_map 	Map
}

func (b Board) GetNeighborIntersections(intersection int) [][2]int {
	thisIntersection := b._map.coordinators[intersection]
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

func (b Board) GetAvailableIntersections(occupied []int) []int {
	intersections := b._map.coordinators
	l := len(occupied)
	for i := 0; i < l; i++ {
		occupiedNeighbors := intersections[occupied[i]].neighbors
		for _, nins := range occupiedNeighbors {
			occupied = append(occupied, nins.index)
		}
	}
	keys := make([]int, 0, len(intersections))
	for k := range intersections {
		if !contains(occupied, k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (b Board) GetTileIndex(intersection int) []int {
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
				tiles[idx].Terrain = ConvertTerrainToCardType(rs[0][2])
			}
		}
	}
	return tiles
}

func contains(s []int, e int) bool {
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

//func GetMap1() Map {
//	gboard := newMap()
//	gboard.build(new(maps.Map1))
//	return gboard
//}