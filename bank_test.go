package main

import (
	"testing"
)

func TestBankGet(t *testing.T) {
	bank := NewBank()
	_, err := bank.Get(0, 1)
	if err != nil || bank.cards[0] != 18 {
		t.Log("expected1 to have 18, failed.")
		t.Fail()
	}
}

func TestBankSet(t *testing.T) {
	bank := NewBank()
	bank.cards[0] = 17
	err := bank.Set(0, 1)
	if err != nil || bank.cards[0] != 19 {
		t.Log("expected to have 19, failed.")
		t.Fail()
	}
}

func TestBankBuyDevCard(t *testing.T) {
	bank := NewBank()
	bank.devCardIndex = 18
	_, err := bank.BuyDevCard()
	if err != nil || bank.devCardIndex != 17 {
		t.Log("expected to have 17, failed.")
		t.Fail()
	}
}

func TestBank1(t *testing.T) {
	bank := NewBank()
	bank.cards[0] = 5

	bank.Begin()
	for i := 0; i < 6; i++ {
		_, err := bank.Get(0, 1)
		if err != nil {
			bank.Rollback()
			break
		}
	}

	if bank.cards[0] != 5 {
		t.Fail()
	}

}
