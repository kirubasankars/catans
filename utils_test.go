package main

import (
	"fmt"
	"testing"
)

func TestRemove(t *testing.T) {
	array := []int{2, 3, 4, 5}
	data := Remove(array, 2)
	if Contains(data, 4) {
		t.Log("Not expected 4")
		t.Fail()
	}

	data = Remove(data, 2)
	if Contains(data, 5) {
		t.Log("Not expected 5")
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	array := []int{1, 3, 5, 7}
	a := Contains(array, 3)
	if a != true {
		t.Log("A expected to be true")
		t.Fail()
	}
}

func TestUnique(t *testing.T) {
	array := []int{1, 4, 7, 7, 9}

	value := Unique(array)
	fmt.Println(len(value))
	if len(value) != 4 {
		t.Log("A expected to be 4")
		t.Fail()
	}
}
