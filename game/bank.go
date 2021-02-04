package game

import (
	"catans/utils"
	"errors"
	"math/rand"
)

type Bank struct {
	cards        [5]int
	devCards     []int

	devCardIndex int
	t            *Bank
}

func (bank *Bank) Begin() {
	bank.t = NewBank()
	bank.t.cards = bank.cards
	bank.t.devCards = bank.devCards
}

func (bank *Bank) Commit() {
	bank.t = nil
}

func (bank *Bank) Rollback() {
	bank.cards = bank.t.cards
	bank.devCards = bank.t.devCards
	bank.t = nil
}

func (bank *Bank) Give(cardType int, count int) (int, error) {
	if bank.cards[cardType] == 0 {
		return 0, nil
	}

	r := bank.cards[cardType] - count
	if r < 0 {
		return -1, errors.New("not enough cards")
	}

	bank.cards[cardType] = r

	return count, nil
}

func (bank *Bank) Return(cardType int, count int) error {
	if bank.cards[cardType] > 19 {
		return errors.New(utils.ErrInvalidOperation)
	}
	bank.cards[cardType] = bank.cards[cardType] + count
	return nil
}

func (bank *Bank) BuyDevCard() (int, error) {
	bank.devCardIndex--
	if bank.devCardIndex <= 0 {
		return -1, errors.New("")
	}
	return bank.devCards[bank.devCardIndex], nil
}

func NewBank() *Bank {
	bank := new(Bank)
	bank.cards = [5]int{19, 19, 19, 19, 19}
	bank.devCards = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 2, 2, 3, 3, 4, 4}
	bank.devCardIndex = len(bank.devCards)
	for i := len(bank.devCards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		bank.devCards[i], bank.devCards[j] = bank.devCards[j], bank.devCards[i]
	}
	return bank
}
