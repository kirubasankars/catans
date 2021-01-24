package game

import "math/rand"

func RandomDice() (int, int) {
	dice1 := rand.Intn(6) + 1
	dice2 := rand.Intn(6) + 1

	return dice1, dice2
}
