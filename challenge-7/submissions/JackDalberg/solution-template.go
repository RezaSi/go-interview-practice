// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"fmt"
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
	return fmt.Sprintf("account error: %s", e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	Balance float64
	Amount float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("insufficient funds: account has $%f.2, attempted to move $%f.2", e.Balance, e.Amount)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf("negative amount: cannot move -$%f.2", -e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	Amount float64
	Limit float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("exceeds limit: limit is %f.2, attempted to more %f.2", e.Limit, e.Amount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" {
	    return nil, &AccountError{Message: "id is empty"}
	}
	if owner == "" {
	    return nil, &AccountError{Message: "owner is empty"}
	}
	if initialBalance < 0.0 {
	    return nil, &NegativeAmountError{Amount: initialBalance}
	}
	if minBalance < 0.0 {
	    return nil, &NegativeAmountError{Amount: minBalance}
	}
	if minBalance > initialBalance {
	    return nil, &InsufficientFundsError{Balance: initialBalance, Amount: minBalance}
	}
	return &BankAccount{
	    ID: id,
	    Owner: owner,
	    Balance: initialBalance,
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
	if a.Balance - amount < a.MinBalance {
	    return &InsufficientFundsError{Amount: amount, Balance: a.Balance}
	}
	a.Balance -= amount
	
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	if err := a.Withdraw(amount); err != nil {
	    return err
	}
	if err := target.Deposit(amount); err != nil {
	    _ = a.Deposit(amount)
	    return err
	}
	return nil
} 