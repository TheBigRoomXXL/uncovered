package something

import (
	"errors"
	"testing"
)

func TestAccount(t *testing.T) {
	var a Account
	if err := a.Deposit(100); err != nil {
		t.Fatal(err)
	}
	if err := a.Withdraw(30); err != nil {
		t.Fatal(err)
	}
	if err := a.Withdraw(1000); !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}
	if a.Balance != 70 {
		t.Fatalf("expected balance 70, got %d", a.Balance)
	}
}
