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
	cards            [5]int
	roads            [][2]int
	settlements      []Settlement
	devCards         []int
	hasLargestArmy   bool
	hasLongestRoad   bool
	knightUsedCount  int
	longestRoadCount int
	ownPort31 bool
	ownPort21 bool
	ports21   [5]int

	allowedSettlementCount        int
	allowedSettlementUpgradeCount int
	allowedRoadsCount			  int

	score     int
}

func (player Player) calculateScore() {
	score := 0
	for _, settlement := range player.settlements {
		score++
		if settlement.Upgraded {
			score++
		}
	}
	for _, devCard := range player.devCards {
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
	for _, card := range player.cards {
		cardCount += card
	}
	return cardCount > limit, cardCount
}

func (player *Player) updateLongestRoad(context GameContext) {
	if len(player.roads) > 4 {
		player.longestRoadCount = context.calculateLongestRoad(*player, []int{})
		for _, otherPlayer := range context.Players {
			if otherPlayer.ID == context.CurrentPlayerID {
				continue
			}

			if player.longestRoadCount > otherPlayer.longestRoadCount {
				player.hasLongestRoad = true
				player.calculateScore()
			}
		}
	}
}

func NewPlayer() *Player {
	p := new(Player)
	p.allowedSettlementCount = 5
	p.allowedSettlementUpgradeCount = 4
	p.allowedRoadsCount = 13
	return p
}
