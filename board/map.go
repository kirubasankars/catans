package board

import "fmt"

type MapSetting interface {
	Tiles() []int
	Connections() [][]int
}

type Map struct {
	nodes        map[int]*Node
	coordinators map[string]*NodeCoordinator
}

func (_map *Map) build(mapSetting MapSetting) {
	tiles := mapSetting.Tiles()
	for tile := range tiles {
		node := NewNode(tile)
		_map.nodes[tile] = node
	}

	connections := mapSetting.Connections()
	for _, item := range connections {
		node := _map.nodes[item[0]]
		neighbor := _map.nodes[item[2]]
		node.neighbors[item[1]] = neighbor
	}

	updateNC := func(node *Node, ncs []*NodeCoordinator) {
		for _, nc := range ncs {
			for i, n := range nc.nodes {
				if n.index == node.index {
					node.coordinators[nc.sides[i]] = nc
				}
			}
		}
	}

	builder := NewCoordinatorBuilder()
	for _, node := range _map.nodes {
		onePointCoordinators := node.findOnePointCoordinators(*builder)
		updateNC(node, onePointCoordinators)
		twoPointsCoordinators := node.findTwoPointsCoordinators(*builder)
		updateNC(node, twoPointsCoordinators)
		threePointsCoordinators := node.findThreePointsCoordinators(*builder)
		updateNC(node, threePointsCoordinators)
	}

	builder.MakeIntersectionsConnected(_map.nodes)

	_map.coordinators = builder.coordinators

	//for i := 0; i < len(_map.nodes); i++ {
	//	fmt.Println(_map.nodes[i])
	//}

	fmt.Println("")
}

func newMap() Map {
	board := new(Map)
	board.nodes = make(map[int]*Node)
	board.coordinators = make(map[string]*NodeCoordinator)
	return *board
}
