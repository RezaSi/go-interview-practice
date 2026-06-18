// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"fmt"
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
	return fmt.Sprintf("Account error: %s", e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	Required  float64
	Available float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("Insufficient funds: need at least %.2f, but only %.2f is available", e.Required, e.Available)

}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf("Negative amount error: %.2f is not a valid amount", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	Amount float64
	Limit  float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("Amount %.2f exceeds transaction limit of %.2f", e.Amount, e.Limit)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" || owner == ""{
	    return nil, &AccountError{"ID and Owner must not be empty"}
	}
	
	
	if initialBalance < 0 {
	    return nil, &NegativeAmountError{Amount: initialBalance}
	}
	if minBalance < 0 {
	    return nil, &NegativeAmountError{Amount: minBalance}
	}
	if initialBalance < minBalance {
	    return nil, &InsufficientFundsError{Required:  minBalance,Available: initialBalance}
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
	// Implement deposit functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{Amount: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Amount: amount, Limit: MaxTransactionAmount}
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
		return &NegativeAmountError{Amount: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Amount: amount, Limit: MaxTransactionAmount}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Required:  a.MinBalance,
			Available: a.Balance - amount,
		}
	}

	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	if a == target {
		return &AccountError{"Cannot transfer to the same account"}
	}

	// Try to withdraw first
	if err := a.Withdraw(amount); err != nil {
		return err
	}

	// If withdraw succeeded, deposit to target
	if err := target.Deposit(amount); err != nil {
		// Rollback the withdrawal if deposit failed
		a.Deposit(amount)
		return err
	}

	return nil
} 