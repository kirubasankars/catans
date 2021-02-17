package main

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Unique(s []int) []int {
	var ns []int
	for _, a := range s {
		if !Contains(ns, a) {
			ns = append(ns, a)
		}
	}
	return ns
}

func Remove(s []int, e int) []int {
	var value = make([]int, len(s)-1)
	i := 0
	for idx, v := range s {
		if idx != e {
			value[i] = v
			i++
		}
	}
	return value
}
