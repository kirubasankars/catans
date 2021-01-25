package game

import (
	"catans/utils"
	"errors"
)

type bank struct {
	cards    [5]int
	devCards []int

	t *bank
}

func (bank *bank) Begin() {
	bank.t = NewBank()
	bank.t.cards = bank.cards
	bank.t.devCards = bank.devCards
}

func (bank *bank) Commit() {
	bank.t = nil
}

func (bank *bank) Rollback() {
	bank.cards = bank.t.cards
	bank.devCards = bank.t.devCards
	bank.t = nil
}

func (bank *bank) Give(cardType int, count int) (int, error) {
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

func (bank *bank) Return(cardType int, count int) error {
	if bank.cards[cardType] > 19 {
		return errors.New(utils.ErrInvalidOperation)
	}
	bank.cards[cardType] = bank.cards[cardType] + count
	return nil
}

func (bank *bank) BuyDevCard() int {
	//r := rand.Intn(len(bank.devCards) + 1)
	//defer utils.Remove(bank.devCards, r)
	//return bank.devCards[r]
	return 0
}

func NewBank() *bank {
	bank := new(bank)
	bank.cards = [5]int{19, 19, 19, 19, 19}
	bank.devCards = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 2, 2, 3, 3, 4, 4}
	return bank
}
