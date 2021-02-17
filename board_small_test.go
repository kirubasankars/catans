package main

import (
	"testing"
)

func TestSmallBoardGetAvailableIntersections(t *testing.T) {
	board := NewBoard(1)
	intersections := board.GetAvailableIntersections([]int{})

	if len(intersections) != 32 {
		t.Log("expected to have 32, failed")
		t.Fail()
	}

	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}

	intersections = board.GetAvailableIntersections([]int{2})
	if len(intersections) != 28 {
		t.Log("expected to have 28, failed")
		t.Fail()
	}
	if Contains(intersections, 1) || Contains(intersections, 9) || Contains(intersections, 3) {
		t.Log("expected not to have  9, 1, 3, failed")
		t.Fail()
	}

	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}

	intersections = board.GetAvailableIntersections([]int{17})
	if len(intersections) != 29 {
		t.Log("expected to have 29, failed")
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
	if len(intersections) != 26 {
		t.Log("expected to have 26, failed")
		t.Fail()
	}
	if Contains(intersections, 17) || Contains(intersections, 16) || Contains(intersections, 4) || Contains(intersections, 2) || Contains(intersections, 14) || Contains(intersections, 4) || Contains(intersections, 3) {
		t.Log("expected not to have 17, 16, 4, 2, 14,3 failed")
		t.Fail()
	}

	if len(Unique(intersections)) != len(intersections) {
		t.Log("expected not to have duplicate, failed")
		t.Fail()
	}
}

func TestSmallBoardGetNeighborIntersections1(t *testing.T) {
	board := NewBoard(1)

	neighbors := board.GetNeighborIntersections1(20)
	if len(neighbors) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	if !(Contains(neighbors, 13) && Contains(neighbors, 21) && Contains(neighbors, 24)) {
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
}

func TestSmallBoardGetNeighborIntersections2(t *testing.T) {
	board := NewBoard(1)

	output := board.GetNeighborIntersections2(20)
	if len(output) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	matched := 0
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

}
func TestSmallGetTileIndices(t *testing.T) {
	board := NewBoard(1)

	indices := board.GetTileIndices(21)
	if len(indices) != 3 {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}
	if !(Contains(indices, 15) && Contains(indices, 20) && Contains(indices, 21)) {
		t.Log("expected to have 3 values, failed")
		t.Fail()
	}

	indices = board.GetTileIndices(4)
	if len(indices) != 2 {
		t.Log("expected to have 1 values, failed")
		t.Fail()
	}
	if !(Contains(indices, 7) && Contains(indices, 13)) {
		t.Log("expected to have 33, failed")
		t.Fail()
	}
}
