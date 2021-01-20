package game

import "errors"

type Bank struct {
	cards    [5]int
	devCards []int

	t 		*Bank
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
		return errors.New("invalid operation")
	}
	bank.cards[cardType] = bank.cards[cardType] + count
	return nil
}

func (bank *Bank) Trade() {

}

func NewBank() *Bank {
	bank := new(Bank)
	bank.cards = [5]int{19,19,19,19,19}
	return bank
}