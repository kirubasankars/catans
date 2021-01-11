package game

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {

	graph1 := new(OccupiedGraph)

	graph1.SetInitialSettlement("1")
	graph1.SetInitialSettlement("6")
	graph1.SetRoad("1", "2")
	graph1.SetRoad("2", "3")
	graph1.SetSettlement("3")
	graph1.SetRoad("6", "7")
	graph1.SetRoad("7", "8")
	fmt.Println(graph1)
}