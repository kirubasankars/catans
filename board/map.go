package board

import (
	"catans/utils"
	"fmt"
)

type TileConnections interface {
	Connections() [][3]int
}

type Map struct {
	nodes        map[int]*Node
	coordinators map[int]*NodeCoordinator
}

func (_map *Map) uniqueTiles(connections [][3]int) []int {
	var nodes []int
	for _, item := range connections {
		if !utils.Contains(nodes, item[0]) {
			nodes = append(nodes, item[0])
		}
	}
	return nodes
}

func (_map *Map) build(tileConnections TileConnections) {
	tiles := _map.uniqueTiles(tileConnections.Connections())
	for tile := range tiles {
		node := NewNode(tile)
		_map.nodes[tile] = node
	}

	connections := tileConnections.Connections()
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

	//_map.coordinators = builder.coordinators

	counter := 0
	for i := 0; i < 19; i++ {
		for j := 0; j < 6; j++ {
			for _, coordinator := range builder.coordinators {
				node := coordinator.nodes[0]
				if node.index == i && coordinator.sides[0] == j {
					coordinator.index = counter
					counter++
				}
			}
		}
	}

	var coordinators = make(map[int]*NodeCoordinator, len(builder.coordinators))
	for _, coordinator := range builder.coordinators {
		coordinators[coordinator.index] = coordinator
	}
	_map.coordinators = coordinators
	//fmt.Println(coordinators)

	//for _, ins := range builder.coordinators {
	//	fmt.Println(ins.index, ins)
	//}

	//for i := 0; i < len(_map.nodes); i++ {
	//	fmt.Println(_map.nodes[i])
	//}

	fmt.Println("")
}

func newMap() Map {
	_map := new(Map)
	_map.nodes = make(map[int]*Node)
	_map.coordinators = make(map[int]*NodeCoordinator)
	return *_map
}
