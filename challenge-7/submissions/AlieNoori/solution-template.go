// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	// Add any other necessary imports
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
type AccountError struct {
	message string
}

func (e *AccountError) Error() string {
	return e.message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	message string
}

func (e *InsufficientFundsError) Error() string {
	return e.message
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	message string
}

func (e *NegativeAmountError) Error() string {
	return e.message
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	message string
}

func (e *ExceedsLimitError) Error() string {
	return e.message
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, &AccountError{"empty id"}
	}

	if owner == "" {
		return nil, &AccountError{"empty Owner"}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{"negative initial balance"}
	}

	if minBalance < 0 {
		return nil, &NegativeAmountError{"negative min balance"}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{"initial balance is less than min balance"}
	}

	bc := &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}

	return bc, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{"negative deposit amount"}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{"deposit amount exceeds transaction limit"}
	}

	a.mu.Lock()
	a.Balance += amount
	a.mu.Unlock()

	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{"negative withdrawal amount"}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{"withdrawal amount exceeds transaction limit"}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{"withdrawal brings the balance below minimum required balance"}
	}

	a.mu.Lock()
	a.Balance -= amount
	a.mu.Unlock()

	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{"negative withdrawal amount"}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{"transfer amount exceeds transaction limit"}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{"transfer brings the balance below minimum required balance"}
	}

	a.mu.Lock()
	a.Balance -= amount
	target.Balance += amount
	a.mu.Unlock()

	return nil
}

