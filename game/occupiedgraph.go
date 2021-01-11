package game

type OccupiedGraph struct {
	name         string
	settlement   bool
	connected    []*OccupiedGraph
}

func (og *OccupiedGraph) SetRoad(intersection1 string, intersection2 string) bool {
	if og.name == intersection1 {
		og1 := new(OccupiedGraph)
		og1.name = intersection2
		og.connected = append(og.connected, og1)
		return true
	} else {
		for _, c := range og.connected {
			r := c.SetRoad(intersection1, intersection2)
			if r == true {
				return true
			}
		}
	}
	return false
}

func (og *OccupiedGraph) SetSettlement(intersection string) bool {
	if og.name == intersection {
		og.settlement = true
		return true
	} else {

		for _, c := range og.connected {
			r := c.SetSettlement(intersection)
			if r == true {
				return true
			}
		}
	}
	return false
}

func (og *OccupiedGraph) SetInitialSettlement(intersection string) {
	og1 := new(OccupiedGraph)
	og1.name = intersection
	og1.settlement = true
	og.connected = append(og.connected, og1)
}

