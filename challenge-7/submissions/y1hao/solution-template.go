// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
)

// BankAccount represents a bank account with balance management and minimum balance requirements.
type BankAccount struct {
	ID         string
	Owner      string
	Balance    float64
	MinBalance float64
	mu         sync.Mutex // For thread safety
}

// Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError string

func (e AccountError) Error() string {
	return string(e)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError string

func (e InsufficientFundsError) Error() string {
	return string(e)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError string

func (e NegativeAmountError) Error() string {
	return string(e)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError string

func (e ExceedsLimitError) Error() string {
	return string(e)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, AccountError("Empty ID")
	}
	if owner == "" {
		return nil, AccountError("Empty owner")
	}
	if initialBalance < 0 {
		return nil, NegativeAmountError("Initial balance is negative")
	}
	if minBalance < 0 {
		return nil, NegativeAmountError("Min balance is negative")
	}
	if initialBalance < minBalance {
		return nil, InsufficientFundsError("Insufficient funds")
	}
	return &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return NegativeAmountError("Amount is negative")
	}

	if amount > MaxTransactionAmount {
		return ExceedsLimitError("The max transaction amount has been exceeded")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return NegativeAmountError("Amount is negative")
	}
	if amount > MaxTransactionAmount {
		return ExceedsLimitError("The max transaction amount has been exceeded")
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return InsufficientFundsError("Insufficient funds")
	}

	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return NegativeAmountError("Amount is negative")
	}
	if amount > MaxTransactionAmount {
		return ExceedsLimitError("The max transaction amount has been exceeded")
	}

	if err := a.Withdraw(amount); err != nil {
		return err
	}

	return target.Deposit(amount)
}
