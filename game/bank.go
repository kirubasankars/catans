package game

import "errors"

type Bank struct {
	cards 			 [5]int
	developmentCards []int
}

func (bank *Bank) Borrow(cardType int, count int) (int, error) {
	if bank.cards[cardType] == 0 {
		return 0, nil
	}

	if bank.cards[cardType] >= count {
		bank.cards[cardType] = bank.cards[cardType] - count
		return count, nil
	}
	if bank.cards[cardType] < 0 {
		r := count - bank.cards[cardType]
		bank.cards[cardType] = 0
		return r, nil
	}
	return 0, nil
}

func (bank *Bank) Return(cardType int, count int) error {
	if bank.cards[cardType] > 19 {
		return errors.New("invalid action")
	}
	bank.cards[cardType] = bank.cards[cardType] + count
	return nil
}

func NewBank() *Bank {
	bank := new(Bank)
	bank.cards = [5]int{19,19,19,19,19}
	return bank
}