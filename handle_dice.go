package main


func (context *GameContext) handleDice(dice int) error {
	context.EventRolled(dice)
	if dice == 7 {
		for _, player := range context.Players {
			if yes, _ := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
				context.scheduleAction(ActionDiscardCards)
				return nil
			}
		}
		context.scheduleAction(ActionPlaceRobber)
		return nil
	}

	bank := context.Bank
	players := context.Players

	bank.Begin()
	defer bank.Commit()

	var cards [][3]int
	for _, player := range players {
		for _, settlement := range player.settlements {
			for idx, token := range settlement.Tokens {
				if token == dice {
					tileIndex := settlement.TileIndex[idx]
					if tileIndex == context.RobberPlacement {
						continue
					}
					terrain := context.Tiles[settlement.TileIndex[idx]][0]
					var count = 0
					var err error = nil
					if settlement.Upgraded {
						count, err = context.Bank.Remove(terrain, 2)
					} else {
						count, err = context.Bank.Remove(terrain, 1)
					}
					if err != nil {
						bank.Rollback()
						return err
					}
					cards = append(cards, [3]int{player.ID, terrain, count})
				}
			}
		}
	}

	for _, card := range cards {
		player := context.Players[card[0]]
		cardType := card[1]
		player.cards[cardType] = player.cards[cardType] + card[2]
		context.EventCardDistributed(player.ID, cardType, card[2])
	}

	return nil
}