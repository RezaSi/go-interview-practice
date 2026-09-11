// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"strings"
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
	Message string
}

func (e *AccountError) Error() string {
	return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Message string
}

func (e *InsufficientFundsError) Error() string {
	return e.Message
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Message string
}

func (e *NegativeAmountError) Error() string {
	return e.Message
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Message string
}

func (e *ExceedsLimitError) Error() string {
	return e.Message
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(owner) == "" {
		return nil, &AccountError{}
	}
	if initialBalance < 0 || minBalance < 0 {
		return nil, &NegativeAmountError{}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{}
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
	if err := a.validateAmount(amount); err != nil {
		return err
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
	if err := a.validateAmount(amount); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{}
	}
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if target == nil {
		return &AccountError{}
	}
	if err := a.validateAmount(amount); err != nil {
		return err
	}

	first, second := a, target
	if target.ID < a.ID {
		first, second = target, a
	}

	first.mu.Lock()
	second.mu.Lock()
	defer first.mu.Unlock()
	defer second.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{}
	}

	a.Balance -= amount
	target.Balance += amount
	return nil
}

func (a *BankAccount) validateAmount(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	return nil
}