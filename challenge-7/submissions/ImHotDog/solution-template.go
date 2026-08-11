// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"errors"
	"fmt"
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
var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
)

// AccountError is a general error type for bank account operations.
type AccountError struct {
	Message string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("AccountError: %s", e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Message string
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("InsufficientFundsError: %s", e.Message)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Message string
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("NegativeAmountError: %s ", e.Message)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Message string
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("ExceedsLimitError: %s ", e.Message)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {

	if id == "" {
		return nil, &AccountError{
			Message: "id is empty",
		}
	}

	if owner == "" {
		return nil, &AccountError{
			Message: "owner is empty",
		}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			Message: "initialBalance must be >= 0",
		}
	}

	if minBalance < 0 {
		return nil, &NegativeAmountError{
			Message: "minBalance must be >= 0",
		}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			Message: "initialBalance must be >= minBalance",
		}
	}

	bankAccount := &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
		mu:         sync.Mutex{},
	}

	return bankAccount, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Message: "transaction amount exceeds maximum allowed limit",
		}
	}

	if amount < 0 {
		return &NegativeAmountError{
			Message: "amount must be positive",
		}
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
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Message: "transaction amount exceeds maximum allowed limit",
		}
	}

	if amount < 0 {
		return &NegativeAmountError{
			Message: "amount must be positive",
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Message: "balance must be >= minBalance after a withdraw",
		}
	}

	a.Balance -= amount

	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Message: "transaction amount exceeds maximum allowed limit",
		}
	}

	if amount < 0 {
		return &NegativeAmountError{
			Message: "amount must be positive",
		}
	}

	first, second := a, target
	if target.ID < a.ID {
		first, second = target, a
	}

	first.mu.Lock()
	defer first.mu.Unlock()

	second.mu.Lock()
	defer second.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Message: "balance must be >= minBalance after a withdraw",
		}
	}

	a.Balance -= amount
	target.Balance += amount

	return nil
}
