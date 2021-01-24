package utils

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Unique(s []int) []int {
	var ns = make([]int, len(s))
	i := 0
	for _, a := range s {
		if !Contains(ns, a) {
			ns[i] = a
			i++
		}
	}
	return ns
}
