package maps

type Map1 struct {
}

func (map1 *Map1) Tiles() []int {
	var nodes []int
	for _, item := range map1.Connections() {
		if !map1.contains(nodes, item[0]) {
			nodes = append(nodes, item[0])
		}
	}
	return nodes
}

func (map1 *Map1) Connections() [][]int {
	return [][]int{
		{0, 1, 1},
		{0, 2, 5},

		{1, 1, 2},
		{1, 2, 6},
		{1, 3, 5},
		{1, 4, 0},

		{2, 1, 3},
		{2, 2, 7},
		{2, 3, 6},
		{2, 4, 1},

		{3, 1, 4},
		{3, 2, 8},
		{3, 3, 7},
		{3, 4, 2},

		{4, 2, 9},
		{4, 3, 8},
		{4, 4, 3},

		{5, 0, 1},
		{5, 1, 6},
		{5, 2, 10},
		{5, 5, 0},

		{6, 0, 2},
		{6, 1, 7},
		{6, 2, 11},
		{6, 3, 10},
		{6, 4, 5},
		{6, 5, 1},

		{7, 0, 3},
		{7, 1, 8},
		{7, 2, 12},
		{7, 3, 11},
		{7, 4, 6},
		{7, 5, 2},

		{8, 0, 4},
		{8, 1, 9},
		{8, 2, 13},
		{8, 3, 12},
		{8, 4, 7},
		{8, 5, 3},

		{9, 2, 14},
		{9, 3, 13},
		{9, 4, 8},
		{9, 5, 4},

		{10, 0, 6},
		{10, 1, 11},
		{10, 2, 15},
		{10, 5, 5},

		{11, 0, 7},
		{11, 1, 12},
		{11, 2, 16},
		{11, 3, 15},
		{11, 4, 10},
		{11, 5, 6},

		{12, 0, 8},
		{12, 1, 13},
		{12, 2, 17},
		{12, 3, 16},
		{12, 4, 11},
		{12, 5, 7},

		{13, 0, 9},
		{13, 1, 14},
		{13, 2, 18},
		{13, 3, 17},
		{13, 4, 12},
		{13, 5, 8},

		{14, 2, 19},
		{14, 3, 18},
		{14, 4, 13},
		{14, 5, 9},

		{15, 0, 11},
		{15, 1, 16},
		{15, 2, 20},
		{15, 5, 10},

		{16, 0, 12},
		{16, 1, 17},
		{16, 2, 21},
		{16, 3, 20},
		{16, 4, 15},
		{16, 5, 11},

		{17, 0, 13},
		{17, 1, 18},
		{17, 2, 22},
		{17, 3, 21},
		{17, 4, 16},
		{17, 5, 12},

		{18, 0, 14},
		{18, 1, 19},
		{18, 2, 23},
		{18, 3, 22},
		{18, 4, 17},
		{18, 5, 13},

		{19, 2, 24},
		{19, 3, 23},
		{19, 4, 18},
		{19, 5, 14},

		{20, 0, 16},
		{20, 1, 21},
		{20, 5, 15},

		{21, 0, 17},
		{21, 1, 22},
		{21, 4, 20},
		{21, 5, 16},

		{22, 0, 18},
		{22, 1, 23},
		{22, 4, 21},
		{22, 5, 17},

		{23, 0, 19},
		{23, 1, 24},
		{23, 4, 22},
		{23, 5, 18},

		{24, 4, 23},
		{24, 5, 19},
	}
}

func (map1 *Map1) contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


