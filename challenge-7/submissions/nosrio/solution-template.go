// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
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

// AccountError is a general error type for bank account operations.
type AccountError struct {
	// Implement this error type
	Reason string
}

func (e *AccountError) Error() string {
	// Implement error message
	return fmt.Sprintf("unable to create account: %s", e.Reason)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	Reason string
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("insufficient founds: %s", e.Reason)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	Reason string
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf("negative ammount: %s", e.Reason)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	Reason string
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("limit exceeded: %s", e.Reason)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" {
		return nil, &AccountError{Reason: "invalid id"}
	}
	if owner == "" {
		return nil, &AccountError{Reason: "invalid owner"}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{Reason: "initial balance can't be negative"}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{Reason: "minimun balance can't be negative"}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{Reason: "initial balance can't be lower than minimum balance"}
	}

	acc := BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
		mu:         sync.Mutex{},
	}
	return &acc, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{Reason: "amount can't be negative"}
	}
	if a.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{Reason: fmt.Sprintf("operation can't be greater than %f", MaxTransactionAmount)}
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
	// Implement withdrawal functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{Reason: "amount can't be negative"}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Reason: fmt.Sprintf("operation can't be greater tha %f", MaxTransactionAmount)}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Reason: fmt.Sprintf("balancer can't be lower than %f", a.MinBalance)}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	err := a.Withdraw(amount)
	if err != nil {
		return err
	}
	err = target.Deposit(amount)
	if err != nil {
		a.Deposit(amount)
		return err
	}
	return nil
}
