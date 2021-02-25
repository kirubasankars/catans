package main

type Player struct {
	ID int
	// 0 - lumber
	// 1 - brick
	// 2 - wool
	// 3 - grain
	// 4 - ore

	// 0 - tree
	// 1 - hill
	// 2 - pasture
	// 3 - field
	// 4 - mountain
	Cards           [5]int
	Roads           [][2]int
	Settlements     []Settlement
	DevCards        []int
	hasLargestArmy  bool
	hasLongestRoad  bool
	KnightUsedCount int

	score     int
	ownPort31 bool
	ownPort21 bool
	ports21   [5]int
}

func (player Player) CalculateScore() {
	score := 0
	for _, settlement := range player.Settlements {
		score++
		if settlement.Upgraded {
			score++
		}
	}
	for _, devCard := range player.DevCards {
		if devCard == 1 {
			score++
		}
	}
	if player.hasLargestArmy {
		score += 2
	}
	if player.hasLongestRoad {
		score += 2
	}
	player.score = score
}

func (player Player) hasMoreCardsThen(limit int) (bool, int) {
	cardCount := 0
	for _, card := range player.Cards {
		cardCount += card
	}
	return cardCount > limit, cardCount
}

func NewPlayer() *Player {
	p := new(Player)
	return p
}
