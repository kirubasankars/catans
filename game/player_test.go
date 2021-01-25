package game

import (
	"fmt"
	"testing"
)

func TestLongestRoad1(t *testing.T) {
	player := new(Player)
	roads := [][2]int{
		{27, 28},
		{28, 29},
		{29, 39},
		{39, 40},
		{38, 39},
		{38, 37},
		{27, 37},
		{38, 48},
		{48, 58},
		//{18, 19},
		//{18, 28},
		//{38, 39},
		//{39, 40},
		//{18, 19},
		//{16, 17},
		//{17, 18},
		//{16, 26},
		//{26, 27},
		//{27, 37},
		//{37, 38},
		//{38, 39},
		//{29, 39},
		//{28, 29},
		//{28, 18},
	}
	player.Roads = append(player.Roads, roads...)

	size := player.calculateLongestRoad(nil)

	if size != 8 {
		t.Error("something wrong")
		t.Fail()
	}
	fmt.Println(size)
}

func TestLongestRoad2(t *testing.T) {
	player := new(Player)
	roads := [][2]int{
		//{38, 48},
		//{39, 40},
		//{18, 19},
		{16, 17},
		{17, 18},
		{16, 26},
		{26, 27},
		{27, 37},
		{37, 38},
		{38, 39},
		{29, 39},
		{28, 29},
		{28, 18},
	}
	player.Roads = append(player.Roads, roads...)

	size := player.calculateLongestRoad([]int{38})

	if size != 10 {
		t.Fail()
	}
	fmt.Println(size)
}
