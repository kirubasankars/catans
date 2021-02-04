package board

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

type hexagon struct {
	index    int
	x        float64
	y        float64
	r        float64
	resource string

	port      bool
	direction float64

	neighbors     []*hexagon
	intersections []*intersection
}

type intersection struct {
	index int
	x     float64
	y     float64
	r     float64

	hasPort      bool
	portResource string

	nodes     []*hexagon
	neighbors []*intersection
}

type grid struct {
	r             float64
	grid          []*hexagon
	intersections []*intersection
}

func (g *grid) makeGrid(h int, w int, tileConfig []string) {
	var r float64 = 50
	var x = r
	var y = r + 20
	var b = 0.27 * r

	var smallRow = true
	var grid = make([]*hexagon, len(tileConfig))
	var idx = 0

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			tc := tileConfig[idx]
			n := hexagon{x: x, y: y, r: r, index: idx, resource: tc}
			if len(tc) == 3 {
				n.port = true
				n.resource = string(tc[1])
				d, _ := strconv.Atoi(string(tc[2]))
				n.direction = float64(d)
			}
			grid[idx] = &n
			x = x + (r * 2) - b
			idx++
		}

		if !smallRow {
			x = r
			smallRow = true
		} else {
			x = (r * 2) - b/2
			smallRow = false
		}
		y = y + (r / 2) + r

	}

	g.grid = grid
}

func (g *grid) makeIntersections() {
	var a = 2 * math.Pi / 6
	var r = g.r
	var ir = r * 0.20

	var id = 0
	var intersections []*intersection
	var intersectionsMap = make(map[string]*intersection)
	for _, h := range g.grid {
		if h.resource == "-" || h.resource == "s" || h.port {
			continue
		}
		var neighbors []*hexagon

		for i := 0.0; i < 6; i++ {
			x := h.x + r*math.Cos((a*i)+11)
			y := h.y + r*math.Sin((a*i)+11)
			var ins *intersection
			var nodes []*hexagon
			var k []string
			var port *hexagon

			for _, h1 := range g.grid {
				if h.resource == "-" || h.resource == "s" {
					continue
				}
				if g.circlesIntersect(x, y, ir, h1.x, h1.y, h1.r) {
					if h1.port {
						port = h1
					}
					nodes = append(nodes, h1)
					k = append(k, strconv.Itoa(h1.index))
				}
			}

			key := strings.Join(k, "#")
			if len(nodes) == 1 {
				key = strconv.Itoa(h.index) + "#" + strconv.Itoa(int(i))
			}
			if v, ok := intersectionsMap[key]; !ok {
				ins = &intersection{index: id, x: x, y: y, r: ir}
				intersectionsMap[key] = ins
				intersections = append(intersections, ins)
				id++
			} else {
				ins = v
			}

			ins.nodes = nodes
			neighbors = append(neighbors, nodes...)
			h.intersections = append(h.intersections, ins)

			if port != nil {
				ins.hasPort = true
				ins.portResource = port.resource
				getNextSide := func(s float64) float64 {
					if s == 5.0 {
						return 0.0
					}
					return s + 1
				}

				px := port.x + r*math.Cos((a*port.direction)+11)
				py := port.y + r*math.Sin((a*port.direction)+11)
				if g.circlesIntersect(ins.x, ins.y, ins.r, px, py, 0) {
					port.intersections = append(port.intersections, ins)
				}

				nd := getNextSide(port.direction)
				px = port.x + r*math.Cos((a*nd)+11)
				py = port.y + r*math.Sin((a*nd)+11)
				if g.circlesIntersect(ins.x, ins.y, ins.r, px, py, 0) {
					port.intersections = append(port.intersections, ins)
				}
			}
		}

		var neighborMap = make(map[int]bool)
		for _, n := range neighbors {
			if n.index == h.index {
				continue
			}
			if _, ok := neighborMap[n.index]; !ok {
				h.neighbors = append(h.neighbors, n)
				neighborMap[n.index] = true
			}
		}
	}

	for _, ins1 := range intersections {
		var neighbors []*intersection
		for _, ins2 := range intersections {
			if ins1.index != ins2.index {
				if g.circlesIntersect(ins1.x, ins1.y, r, ins2.x, ins2.y, ins2.r) {
					neighbors = append(neighbors, ins2)
				}
			}
		}
		ins1.neighbors = neighbors
	}
	g.intersections = intersections
}

func (g *grid) circlesIntersect(x1, y1, r1, x2, y2, r2 float64) bool {
	return (x1-x2)*(x1-x2)+(y1-y2)*(y1-y2) <= (r1+r2)*(r1+r2)
}

func (g *grid) parse(m string) ([]string, int, int) {
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

func (g *grid) Build(tiles string) {
	g.r = 50
	config, height, width := g.parse(tiles)
	g.makeGrid(height, width, config)
	g.makeIntersections()
}
