package main

import "math/rand"

func (context *GameContext) giveInitialFreeCards() error {
	context.Bank.Begin()
	for _, player := range context.Players {
		r := rand.Intn(2)
		indices := player.settlements[r].TileIndex

		giveCard := func(idx int) {
			cardType := context.Tiles[indices[idx]][0]
			context.Bank.Remove(cardType, 1)
			player.cards[cardType]++
			context.EventCardDistributed(player.ID, cardType, 1)
		}

		giveCard(0)
		if len(indices) > 1 {
			giveCard(1)
		}
		if len(indices) > 2 {
			giveCard(2)
		}
	}
	context.Bank.Commit()
	return nil
}

