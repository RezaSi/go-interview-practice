// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"unicode"
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
	// Implement this error type
	S string
}

func (e *AccountError) Error() string {
	// Implement error message
	return e.S
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	S string
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return e.S
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	S string
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return e.S
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	S string
	// Implement this error type
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return e.S
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if !valid(id) || !valid(owner) {
		return nil, &AccountError{S: "Wrong ID or owner"}
	}

	if initialBalance < 0 || minBalance < 0 {
		return nil, &NegativeAmountError{S: "Negative amount error"}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{S: "IFE"}
	}

	return &BankAccount{
			ID:         id,
			Owner:      owner,
			Balance:    initialBalance,
			MinBalance: minBalance,
		},
		nil
}

func valid(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{S: "Negative amount error"}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{S: "ELE"}
	}
	if a == nil {
		return &AccountError{S: "nil"}
	}

	defer a.mu.Unlock()
	a.mu.Lock()
	a.Balance += amount

	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Implement withdrawal functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{S: "Negative amount error"}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{S: "ELE"}
	}
	if a == nil {
		return &AccountError{S: "nil"}
	}

	defer a.mu.Unlock()
	a.mu.Lock()
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{S: "IFE"}
	}
	a.Balance -= amount

	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if a == nil || target == nil {
		return &AccountError{S: "nil"}
	}
	if amount < 0 {
		return &NegativeAmountError{S: "Negative amount error"}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{S: "ELE"}
	}
	// Implement transfer functionality with proper error handling
	defer a.mu.Unlock()
	a.mu.Lock()
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{S: "IFE"}
	}
	a.Balance -= amount

	defer target.mu.Unlock()
	target.mu.Lock()

	target.Balance += amount
	return nil
}