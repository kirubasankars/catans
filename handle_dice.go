package main


func (context *GameContext) handleDice(dice int) error {
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
		for _, settlement := range player.Settlements {
			for idx, token := range settlement.Tokens {
				if token == dice {
					tileIndex := settlement.Indices[idx]
					if tileIndex == context.RobberPlacement {
						continue
					}
					terrain := context.Tiles[settlement.Indices[idx]][0]
					var count = 0
					var err error = nil
					if settlement.Upgraded {
						count, err = context.Bank.Get(terrain, 2)
					} else {
						count, err = context.Bank.Get(terrain, 1)
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
		player.Cards[cardType] = player.Cards[cardType] + card[2]
	}

	return nil
}