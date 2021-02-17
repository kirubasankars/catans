package main

import "math/rand"

func (context *GameContext) randomDiscardCards() {
	for _, player := range context.Players {
		if yes, cardCount := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
			numberOfCardsRemove := cardCount / 2

			for {
				cardId := rand.Intn(5)
				card := player.Cards[cardId]
				if card == 0 {
					if numberOfCardsRemove <= 0 {
						break
					}
					continue
				}

				r := rand.Intn(card) + 1
				if numberOfCardsRemove-r >= 0 {
					player.Cards[cardId] = card - r
					numberOfCardsRemove -= r
				} else {
					player.Cards[cardId] = card - numberOfCardsRemove
					numberOfCardsRemove = 0
				}

				if numberOfCardsRemove <= 0 {
					break
				}
			}
		}
	}
}
