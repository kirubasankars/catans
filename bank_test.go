package main

import "testing"

func TestBankGive(t *testing.T) {
	bank := NewBank()
	_, err := bank.Give(0, 1)
	if err != nil || bank.cards[0] != 18 {
		t.Log("expected to have 18, failed.")
		t.Fail()
	}
}
