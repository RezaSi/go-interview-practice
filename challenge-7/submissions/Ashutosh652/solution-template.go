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
	// Implement this error type
	Message string
}

func (e *AccountError) Error() string {
	// Implement error message
	return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	Message string
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return e.Message
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	Message string
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return e.Message
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	Message string
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return e.Message
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" {
		return nil, &AccountError{Message: "Account ID cannot be empty."}
	}
	if owner == "" {
		return nil, &AccountError{Message: "Account owner cannot be empty."}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{Message: "Minimum balance cannot be negative."}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{Message: "Initial balance cannot be negative."}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{Message: "Balance is less than minimum allowed balance."}
	}
	return &BankAccount{ID: id, Owner: owner, Balance: initialBalance, MinBalance: minBalance}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{Message: "Deposit amount cannot be negative."}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Message: "Deposit amount exceeds the maximum transaction limit."}
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
		return &NegativeAmountError{Message: "Withdrawal amount cannot be negative."}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Message: "Withdrawal amount exceeds the maximum transaction limit."}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Message: "Insufficient funds for this withdrawal."}
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
	if a == target {
		return &AccountError{Message: "Cannot transfer to the same account."}
	}
	if amount < 0 {
		return &NegativeAmountError{Message: "Transfer amount cannot be negative."}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Message: "Transfer amount exceeds the maximum transaction limit."}
	}

	a.mu.Lock()
	target.mu.Lock()
	defer a.mu.Unlock()
	defer target.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Message: "Insufficient funds for this transfer."}
	}
	a.Balance -= amount
	target.Balance += amount

	return nil
}
