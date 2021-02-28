package main

import "math/rand"

func (context *GameContext) giveInitialFreeCards() error {
	context.Bank.Begin()
	for _, player := range context.Players {
		r := rand.Intn(2)
		indices := player.Settlements[r].Indices
		cardType := context.Tiles[indices[0]][0]
		context.Bank.Get(cardType, 1)
		player.Cards[cardType]++
		context.EventCardDistributed(player.ID, cardType, 1)

		if len(indices) > 1 {
			cardType := context.Tiles[indices[1]][0]
			context.Bank.Get(cardType, 1)
			player.Cards[cardType]++

			context.EventCardDistributed(player.ID, cardType, 1)
		}
	}
	context.Bank.Commit()
	return nil
}

