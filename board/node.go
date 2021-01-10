package board

import "fmt"

type Node struct {
	index        int
	neighbors    map[int]*Node
	coordinators [6]*NodeCoordinator
}

func (node Node) findOnePointCoordinators(builder CoordinatorBuilder) []*NodeCoordinator {
	var coords []*NodeCoordinator
	for currentSide := 0; currentSide <= 5; currentSide++ {
		if _, ok := node.neighbors[currentSide]; !ok {
			nextSide := node.getNextSide(currentSide)
			if _, ok := node.neighbors[nextSide]; !ok {
				nc := builder.GetNodeCoordinator([]*Node{&node}, []int{nextSide})
				coords = append(coords, nc)
			}
		}
	}
	return coords
}

func (node Node) findTwoPointsCoordinators(builder CoordinatorBuilder) []*NodeCoordinator {
	currentNode := &node
	var coords []*NodeCoordinator
	for side := 0; side <= 5; side++ {
		if _, ok := currentNode.neighbors[side]; !ok {
			nextSide := currentNode.getNextSide(side)
			if neighbor, ok := currentNode.neighbors[nextSide]; ok {
				for neighborSide, nNeighbor := range neighbor.neighbors {
					if nNeighbor.index == currentNode.index {

						currentNodeSide := nextSide
						neighborNodeSide := currentNode.getNextSide(neighborSide)

						var nc *NodeCoordinator
						if currentNode.index < neighbor.index {
							nc = builder.GetNodeCoordinator([]*Node{currentNode, neighbor}, []int{currentNodeSide, neighborNodeSide})
						} else if neighbor.index < currentNode.index {
							nc = builder.GetNodeCoordinator([]*Node{neighbor, currentNode}, []int{neighborNodeSide, currentNodeSide})
						}
						coords = append(coords, nc)
					}
				}
			}

		}
	}

	for side := 0; side <= 5; side++ {
		if neighbor, ok := currentNode.neighbors[side]; ok {
			nextSide := currentNode.getNextSide(side)
			if _, ok := currentNode.neighbors[nextSide]; !ok {
				for side, nNeighbor := range neighbor.neighbors {
					if nNeighbor.index == currentNode.index {
						currentNodeSide := nextSide
						neighborNodeSide := side

						var nc *NodeCoordinator
						if currentNode.index < neighbor.index {
							nc = builder.GetNodeCoordinator([]*Node{currentNode, neighbor}, []int{currentNodeSide, neighborNodeSide})
						} else if neighbor.index < currentNode.index {
							nc = builder.GetNodeCoordinator([]*Node{neighbor, currentNode}, []int{neighborNodeSide, currentNodeSide})
						}
						coords = append(coords, nc)
					}
				}
			}
		}
	}

	return coords
}

func (node Node) findThreePointsCoordinators(builder CoordinatorBuilder) []*NodeCoordinator {
	var coords []*NodeCoordinator
	for _, neighbor1 := range node.neighbors {
		for _, neighbor2 := range neighbor1.neighbors {
			if neighbor2.index == node.index {
				continue
			}
			for _, neighbor3 := range neighbor2.neighbors {
				if neighbor3.index == node.index {
					var (
						xNode *Node
						yNode *Node
						zNode *Node
						xSide int
						ySide int
						zSide int
					)

					if neighbor1.index < neighbor2.index && neighbor1.index < neighbor3.index {
						xNode = neighbor1
					} else if neighbor2.index < neighbor1.index && neighbor2.index < neighbor3.index {
						xNode = neighbor2
					} else if neighbor3.index < neighbor1.index && neighbor3.index < neighbor2.index {
						xNode = neighbor3
					}

					if neighbor1.index > neighbor2.index && neighbor1.index > neighbor3.index {
						zNode = neighbor1
					} else if neighbor2.index > neighbor1.index && neighbor2.index > neighbor3.index {
						zNode = neighbor2
					} else if neighbor3.index > neighbor1.index && neighbor3.index > neighbor2.index {
						zNode = neighbor3
					}

					if xNode.index < neighbor1.index && neighbor1.index < zNode.index {
						yNode = neighbor1
					} else if xNode.index < neighbor2.index && neighbor2.index < zNode.index {
						yNode = neighbor2
					} else if xNode.index < neighbor3.index && neighbor3.index < zNode.index {
						yNode = neighbor3
					}

					for xs, n := range xNode.neighbors {
						if n.index == yNode.index {
							xSide = xs
						}
					}
					for ys, n := range yNode.neighbors {
						if n.index == zNode.index {
							ySide = ys
						}
					}
					for zs, n := range zNode.neighbors {
						if n.index == xNode.index {
							zSide = zs
						}
					}

					nc := builder.GetNodeCoordinator([]*Node{xNode, yNode, zNode}, []int{xSide, ySide, zSide})
					coords = append(coords, nc)
				}
			}
		}
	}

	var coords1 []*NodeCoordinator
	dedup := make(map[string]*NodeCoordinator)
	for _, nc := range coords {
		key := fmt.Sprint(nc)
		if _, ok := dedup[key]; !ok {
			dedup[key] = nc
			coords1 = append(coords1, nc)
		}
	}

	return coords1
}

func (node Node) getNextSide(side int) int {
	if side >= 5 {
		return 0
	}
	return side + 1
}

func (node Node) getPrevSide(side int) int {
	if side <= 0 {
		return 5
	}
	return side - 1
}

func (node Node) String() string {
	return fmt.Sprint(node.index, node.coordinators)
}

func NewNode(index int) *Node {
	n := new(Node)
	n.index = index
	n.neighbors = make(map[int]*Node)
	return n
}
