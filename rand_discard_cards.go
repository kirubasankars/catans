package main

import "math/rand"

func (context *GameContext) randomDiscardCards() {
	for _, player := range context.Players {
		if yes, cardCount := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
			numberOfCardsRemove := cardCount / 2

			for {
				cardType := rand.Intn(5)
				cardCount := player.cards[cardType]
				if cardCount == 0 {
					if numberOfCardsRemove <= 0 {
						break
					}
					continue
				}

				randCardCount2Remove := rand.Intn(cardCount) + 1
				if numberOfCardsRemove-randCardCount2Remove >= 0 {
					player.cards[cardType] = cardCount - randCardCount2Remove
					numberOfCardsRemove -= randCardCount2Remove
				} else {
					player.cards[cardType] = cardCount - numberOfCardsRemove
					numberOfCardsRemove = 0
				}

				if numberOfCardsRemove <= 0 {
					break
				}
			}
		}
	}
}
