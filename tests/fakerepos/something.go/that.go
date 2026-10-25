package something

import "errors"

// Account is a toy bank account.
type Account struct {
	Balance int
}

var ErrInsufficientFunds = errors.New("insufficient funds")

// Deposit is partially covered: negative amounts are never tested.
func (a *Account) Deposit(amount int) error {
	if amount <= 0 {
		return errors.New("deposit must be positive")
	}
	a.Balance += amount
	return nil
}

// Withdraw is fully covered.
func (a *Account) Withdraw(amount int) error {
	if amount > a.Balance {
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	return nil
}

// Transfer is never called by the tests.
func Transfer(from, to *Account, amount int) error {
	if err := from.Withdraw(amount); err != nil {
		return err
	}
	if err := to.Deposit(amount); err != nil {
		from.Balance += amount
		return err
	}
	return nil
}
