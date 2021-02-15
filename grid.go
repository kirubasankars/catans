package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Hexagon struct {
	index   int
	x       float64
	y       float64
	r       float64
	terrain string
	token   int
	port    *HexagonPort

	//neighbors     []*Hexagon
	//intersections []*Intersection
}

type HexagonPort struct {
	direction float64
	resource  string
}

type Intersection struct {
	index int
	x     float64
	y     float64
	r     float64

	port *Hexagon

	nodes     []*Hexagon
	neighbors []*Intersection
}

type MapConfig interface {
	GetTileConfig() string
	GetChits() []int
}

type Grid struct {
	nodes         []*Hexagon
	intersections []*Intersection
}

func (grid *Grid) makeGrid(h int, w int, tileConfig []string, tokens []int) {
	var r float64 = 50
	var x = r
	var y = r + 20
	var b = 0.27 * r

	var alternativeRow = true
	var nodes = make([]*Hexagon, len(tileConfig))
	var idx = 0
	var tokenIdx = 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			tc := tileConfig[idx]
			n := Hexagon{x: x, y: y, r: r, index: idx, terrain: tc}
			if len(tc) == 3 {
				d, _ := strconv.Atoi(string(tc[2]))
				n.port = &HexagonPort{direction: float64(d), resource: string(tc[1])}
				n.terrain = string(tc[1])
			}

			if !(n.terrain == "d" || n.terrain == "-" || n.terrain == "s" || n.port != nil) { //ignore any thing is not tokens
				n.token = tokens[tokenIdx]
				tokenIdx++
			}

			nodes[idx] = &n
			x = x + (r * 2) - b
			idx++
		}

		if !alternativeRow {
			x = r
			alternativeRow = true
		} else {
			x = (r * 2) - b/2
			alternativeRow = false
		}
		y = y + (r / 2) + r
	}

	if tokenIdx < len(tokens) {
		panic(ErrInvalidOperation)
	}

	grid.nodes = nodes
}

func (grid *Grid) makeIntersections() {
	var a = 2 * math.Pi / 6
	var r = 50.0
	var ir = r * 0.20

	var id = 0
	var intersections []*Intersection
	var intersectionsMap = make(map[string]*Intersection)

	for _, node := range grid.nodes {
		if node.terrain == "-" || node.terrain == "s" || node.port != nil {
			continue
		}
		//var neighbors []*Hexagon

		for i := 0.0; i < 6; i++ {
			x := node.x + r*math.Cos((a*i)+11)
			y := node.y + r*math.Sin((a*i)+11)

			var nodes []*Hexagon
			var ins *Intersection
			var k []string

			for _, h1 := range grid.nodes {
				if grid.circlesIntersect(x, y, ir, h1.x, h1.y, h1.r) {
					nodes = append(nodes, h1)
					k = append(k, strconv.Itoa(h1.index))
				}
			}

			key := strings.Join(k, "#")
			if v, ok := intersectionsMap[key]; !ok {
				ins = &Intersection{index: id, x: x, y: y, r: ir}
				intersectionsMap[key] = ins
				intersections = append(intersections, ins)
				id++
			} else {
				ins = v
			}

			ins.nodes = nodes
			//neighbors = append(neighbors, nodes...)
			//node.intersections = append(node.intersections, ins)
		}

		//ignore myself from neighbors
		//for _, n := range neighbors {
		//	if n.index == node.index {
		//		continue
		//	}
		//	node.neighbors = append(node.neighbors, n)
		//}
	}

	getNextSide := func(s float64) float64 {
		if s == 5.0 {
			return 0.0
		}
		return s + 1
	}

	for _, node := range grid.nodes {
		if node.port == nil {
			continue
		}
		var portIntersections []*Intersection
		for _, ins := range intersections {
			if grid.circlesIntersect(node.x, node.y, node.r, ins.x, ins.y, ins.r) {
				portIntersections = append(portIntersections, ins)
			}
		}

		for _, ins := range portIntersections {
			var px = node.x + r*math.Cos((a*node.port.direction)+11)
			var py = node.y + r*math.Sin((a*node.port.direction)+11)
			if grid.circlesIntersect(ins.x, ins.y, ins.r, px, py, ins.r) {
				ins.port = node
			}

			var ns = getNextSide(node.port.direction)
			px = node.x + r*math.Cos((a*ns)+11)
			py = node.y + r*math.Sin((a*ns)+11)
			if grid.circlesIntersect(ins.x, ins.y, ins.r, px, py, ins.r) {
				ins.port = node
			}
		}
	}

	//find neighbors intersections
	for _, ins1 := range intersections {
		var neighbors []*Intersection
		for _, ins2 := range intersections {
			if ins1.index != ins2.index {
				if grid.circlesIntersect(ins1.x, ins1.y, r, ins2.x, ins2.y, ins2.r) {
					neighbors = append(neighbors, ins2)
				}
			}
		}
		ins1.neighbors = neighbors
	}

	grid.intersections = intersections
}

func (grid *Grid) circlesIntersect(x1, y1, r1, x2, y2, r2 float64) bool {
	return (x1-x2)*(x1-x2)+(y1-y2)*(y1-y2) <= (r1+r2)*(r1+r2)
}

func (grid *Grid) parseTiles(m string) ([]string, int, int) {
	rx, _ := regexp.Compile("[[:blank:]]")
	o := rx.ReplaceAll([]byte(m), []byte(""))
	var output [][]string
	segments := strings.Split(string(o), "\n")

	for _, seg := range segments {
		row := strings.Split(seg, ",")
		if len(row) > 1 {
			var newRow []string
			for _, v := range row {
				if len(v) == 0 {
					continue
				}
				newRow = append(newRow, v)
			}
			output = append(output, newRow)
		}
	}

	var tiles []string
	for _, r := range output {
		for _, s := range r {
			tiles = append(tiles, s)
		}
	}
	h := len(output)
	w := len(output[0])

	return tiles, h, w
}

func (grid *Grid) Build(config MapConfig) {
	tilesConfig, height, width := grid.parseTiles(config.GetTileConfig())
	chits := config.GetChits()
	grid.makeGrid(height, width, tilesConfig, chits)
	grid.makeIntersections()
}
