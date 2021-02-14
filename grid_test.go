package main

import "testing"

type testMap1 struct{}

func (tmap testMap1) GetTileConfig() string {
	return `
	-,-,-
	-,t,-
	-,-,-`
}

func (tmap testMap1) GetChits() []int {
	return []int{6}
}

func TestOneTileGrid(t *testing.T) {
	grid := Grid{}
	grid.Build(testMap1{})

	if len(grid.intersections) != 6 {
		t.Fail()
	}
}

type testMap2 struct{}

func (tmap testMap2) GetTileConfig() string {
	return `
	-,-,-,-
	-,t,p,-
	-,-,-,-`
}

func (tmap testMap2) GetChits() []int {
	return []int{6, 8}
}

func TestTwoTileGrid(t *testing.T) {
	grid := Grid{}
	grid.Build(testMap2{})

	if len(grid.intersections) != 10 {
		t.Fail()
	}
}

type testMap3 struct{}

func (tmap testMap3) GetTileConfig() string {
	return `
	-,-,-,-
	-,t,-,-
	-,t,p,-
	-,-,-,-`
}

func (tmap testMap3) GetChits() []int {
	return []int{6, 8, 9}
}

func TestThreeTileGrid(t *testing.T) {
	grid := Grid{}
	grid.Build(testMap3{})

	if len(grid.intersections) != 13 {
		t.Fail()
	}
}

func TestSmallMap(t *testing.T) {
	grid := Grid{}
	grid.Build(SmallMap{})

	if len(grid.intersections) != 32 {
		t.Fail()
	}
}
