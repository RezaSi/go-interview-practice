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
	msg string
}

func (e *AccountError) Error() string {
	// Implement error message
	return e.msg
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	Balance    float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("Your balance is %f , it can't be less than %f ", e.Balance, e.MinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf("excepted positive amount , got %f", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	Amount float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("%f exceeds defined limit of operation", e.Amount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" || owner == "" {
		return nil, &AccountError{
			msg: "ID and Owner fields required",
		}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{
			Amount: minBalance,
		}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			Amount: initialBalance,
		}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			Balance:    initialBalance,
			MinBalance: minBalance,
		}
	}
	bankAccount := &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}
	return bankAccount, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	// Implement deposit functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{
			Amount: amount,
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
		}
	}
	a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount < 0 {
		return &NegativeAmountError{
			Amount: amount,
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
		}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Balance:    a.Balance,
			MinBalance: a.MinBalance,
		}
	}
	// Implement withdrawal functionality with proper error handling
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount < 0 {
		return &NegativeAmountError{
			Amount: amount,
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
		}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Balance:    a.Balance,
			MinBalance: a.MinBalance,
		}
	}
	a.Balance -= amount
	target.Balance += amount
	return nil
}
