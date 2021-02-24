package main

import "math/rand"

func (context *GameContext) randomDiscardCards() {
	for _, player := range context.Players {
		if yes, cardCount := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
			numberOfCardsRemove := cardCount / 2

			for {
				cardType := rand.Intn(5)
				cardCount := player.Cards[cardType]
				if cardCount == 0 {
					if numberOfCardsRemove <= 0 {
						break
					}
					continue
				}

				randCardCount2Remove := rand.Intn(cardCount) + 1
				if numberOfCardsRemove-randCardCount2Remove >= 0 {
					player.Cards[cardType] = cardCount - randCardCount2Remove
					numberOfCardsRemove -= randCardCount2Remove
				} else {
					player.Cards[cardType] = cardCount - numberOfCardsRemove
					numberOfCardsRemove = 0
				}

				if numberOfCardsRemove <= 0 {
					break
				}
			}
		}
	}
}
