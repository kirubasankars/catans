package board

import (
	"catans/maps"
	"catans/utils"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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

type Hexagon struct {
	NodeID        int               `json:"node_id"`
	Token         int				`json:"token"`
	Terrain       string			`json:"terrain"`
	Neighbors     []HexagonNeighbor `json:"neighbors,omitempty"`
	Intersections []HexagonIntersection `json:"intersections,omitempty"`
}

type HexagonIntersection struct {
	IntersectionID int `json:"intersection_id"`
	Side   		   int `json:"side"`
}

type HexagonNeighbor struct {
	NodeID int `json:"node_id"`
	Side   int `json:"side"`
}

func (board Board) GetUINodes() {
	var hexMap []Hexagon
	var visitedNodes = make(map[int]bool)
	var visitedintersections = make(map[int]bool)

	tiles := generateTiles("10MO,2PA,9FO,12FI,6HI,4PA,10HI,9FI,11FO,0DE,3FO,8MO,8FO,3MO,4FI,5PA,5HI,6FI,11PA")

	for i := 0; i < len(board._map.nodes); i++ {
		node := board._map.nodes[i]
		tile := tiles[i]
		var hex Hexagon
		hex.Token = tile.Token
		hex.Terrain = convertCardTypeToName(tile.Terrain)
		hex.NodeID = node.index
		for side, n := range node.neighbors {
			if _, ok := visitedNodes[n.index]; !ok {
				hex.Neighbors = append(hex.Neighbors, HexagonNeighbor{n.index, side})
				visitedNodes[n.index] = true
			}
			visitedNodes[node.index] = true
		}

		for side, i := range node.coordinators {
			if _, ok := visitedintersections[i.index]; !ok {
				hex.Intersections = append(hex.Intersections, HexagonIntersection{i.index, side})
				visitedintersections[i.index] = true
			}
		}

		hexMap = append(hexMap, hex)
	}

	s, _ := json.Marshal(hexMap)

	fmt.Println(string(s))
}

type tile struct {
	Token   int
	Terrain int
}

func convertTerrainToCardType(terrain string) int {
	switch terrain {
	case "FO":
		return 0
	case "HI":
		return 1
	case "PA":
		return 2
	case "FI":
		return 3
	case "MO":
		return 4
	}
	return -1
}

func convertCardTypeToName(cardType int) string {
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

func generateTiles(tileSettings string) []tile {
	r := regexp.MustCompile(`(?P<Token>\d+)(?P<Terrain>\w+)?`)
	segs := strings.Split(tileSettings, ",")
	tiles := make([]tile, len(segs))
	for idx, seg := range segs {
		rs := r.FindAllStringSubmatch(seg, -1)
		if len(rs) > 0 {
			tiles[idx].Token, _ = strconv.Atoi(rs[0][1])
			if len(rs[0]) > 1 {
				tiles[idx].Terrain = convertTerrainToCardType(rs[0][2])
			}
		}
	}
	return tiles
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
