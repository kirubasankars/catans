package main

import (
	"testing"
)

func TestBoardGetAvailableIntersections(t *testing.T) {
	board := NewBoard(0)
	intersections := board.GetAvailableIntersections([]int{})

	if len(intersections) != 54 {
		t.Log("expected to have 54, failed")
		t.Fail()
	}

	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}

	intersections = board.GetAvailableIntersections([]int{2})
	if len(intersections) != 50{
		t.Log("expected to have 50, failed")
		t.Fail()
	}
	if Contains(intersections, 2) || Contains(intersections, 9) || Contains(intersections, 1) || Contains(intersections, 3) {
		t.Log("expected not to have 2, 9, 1, 3, failed")
		t.Fail()
	}

	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}

	intersections = board.GetAvailableIntersections([]int{17})
	if len(intersections) != 51 {
		t.Log("expected to have 52, failed")
		t.Fail()
	}
	if Contains(intersections, 17) || Contains(intersections, 16) || Contains(intersections, 4) {
		t.Log("expected not to have 17, 16, 4, failed")
		t.Fail()
	}
	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}

	intersections = board.GetAvailableIntersections([]int{17, 3})
	if len(intersections) != 48 {
		t.Log("expected to have 48, failed")
		t.Fail()
	}
	if Contains(intersections, 17) || Contains(intersections, 16) || Contains(intersections, 4) || Contains(intersections, 2) || Contains(intersections, 14) {
		t.Log("expected not to have 17, 16, 4, 2, 14 failed")
		t.Fail()
	}

	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}
}

func TestBoardGetNeighborIntersections1(t *testing.T) {
	board := NewBoard(0)

	neighbors := board.GetNeighborIntersections1(20)
	if len(neighbors) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	if !(Contains(neighbors, 13) && Contains(neighbors, 21) && Contains(neighbors, 24)) {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}

	neighbors = board.GetNeighborIntersections1(32)
	if len(neighbors) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	if !(Contains(neighbors, 29) && Contains(neighbors, 31) && Contains(neighbors, 41)) {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}

	neighbors = board.GetNeighborIntersections1(10)
	if len(neighbors) != 2 {
		t.Log("expected to have 2 values, failed")
		t.Fail()
	}
	if !(Contains(neighbors, 7) && Contains(neighbors, 11)) {
		t.Log("expected to have 2 values, failed")
		t.Fail()
	}

	neighbors = board.GetNeighborIntersections1(46)
	if len(neighbors) != 3 {
		t.Log("expected to have 2 values, failed")
		t.Fail()
	}
	if !(Contains(neighbors, 52) && Contains(neighbors, 45) && Contains(neighbors, 43)) {
		t.Log("expected to have 2 values, failed")
		t.Fail()
	}

}

func TestBoardGetNeighborIntersections2(t *testing.T) {
	board := NewBoard(0)

	output := board.GetNeighborIntersections2(20)
	if len(output) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	matched  := 0
	for _, neighbors := range output {
		if neighbors[0] == 13 || neighbors[1] == 13 {
			matched++
		}
		if neighbors[0] == 21 || neighbors[1] == 21 {
			matched++
		}
		if neighbors[0] == 24 || neighbors[1] == 24 {
			matched++
		}
	}

	if matched != 3 {
		t.Fail()
	}

	output = board.GetNeighborIntersections2(32)
	if len(output) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}

	matched  = 0
	for _, neighbors := range output {
		if neighbors[0] == 29 || neighbors[1] == 29 {
			matched++
		}
		if neighbors[0] == 31 || neighbors[1] == 31 {
			matched++
		}
		if neighbors[0] == 41 || neighbors[1] == 41 {
			matched++
		}
	}

	if matched != 3 {
		t.Fail()
	}
}

func TestGetTileIndices(t *testing.T) {
	board := NewBoard(0)

	indices := board.GetTileIndices(21)
	if len(indices) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	if !(Contains(indices, 24) && Contains(indices, 25) && Contains(indices, 18)) {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}

	indices = board.GetTileIndices(45)
	if len(indices) != 2 {
		t.Log("expected to have 1 values, failed")
		t.Fail()
	}
	if !(Contains(indices, 33) && Contains(indices, 34)) {
		t.Log("expected to have 33, failed")
		t.Fail()
	}
}