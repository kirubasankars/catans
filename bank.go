package main

import (
	"errors"
	"math/rand"
)

type Bank struct {
	// 0 - lumber
	// 1 - brick
	// 2 - wool
	// 3 - grain
	// 4 - ore
	cards    [5]int
	devCards []int

	devCardIndex int
	t            *Bank
}

func (bank *Bank) Begin() {
	bank.t = &Bank{}
	bank.t.cards = bank.cards
	bank.t.devCards = bank.devCards
	bank.t.devCardIndex = bank.devCardIndex
}

func (bank *Bank) Commit() {
	bank.t = nil
}

func (bank *Bank) Rollback() {
	bank.cards = bank.t.cards
	bank.devCardIndex = bank.t.devCardIndex
	bank.devCards = bank.t.devCards
	bank.t = nil
}

func (bank *Bank) Get(cardType int, count int) (int, error) {
	if bank.cards[cardType] == 0 {
		return 0, errors.New("not enough cards")
	}

	r := bank.cards[cardType] - count
	if r < 0 {
		return -1, errors.New("not enough cards")
	}

	bank.cards[cardType] = r

	return count, nil
}

func (bank *Bank) Set(cardType int, count int) error {
	if bank.cards[cardType]+count > 19 {
		return errors.New(ErrInvalidOperation)
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
	bank.devCards = []int{DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardKnight, DevCardVPPoint, DevCardVPPoint, DevCardVPPoint, DevCardVPPoint, DevCardVPPoint, DevCardMonopoly, DevCardMonopoly, DevCard2Road, DevCard2Road, DevCard2Resource, DevCard2Resource}
	bank.devCardIndex = len(bank.devCards)
	for i := len(bank.devCards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		bank.devCards[i], bank.devCards[j] = bank.devCards[j], bank.devCards[i]
	}
	return bank
}
