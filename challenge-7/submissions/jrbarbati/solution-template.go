// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
    "fmt"
    "math"
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
type AccountError struct {
	message string
}

func (e AccountError) Error() string {
	return e.message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	minBalance float64
}

func (e InsufficientFundsError) Error() string {
	return fmt.Sprintf("unable to perform withdrawal as balance would go below minimum allowed (%v)", e.minBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
}

func (e NegativeAmountError) Error() string {
	return "amount cannot be negative"
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
}

func (e ExceedsLimitError) Error() string {
	return "exceeds maximum transaction limit"
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
    if id == "" || owner == "" {
        return nil, AccountError{message: "Account needs both an ID and an Owner"}
    }
    
    if !isFinite(initialBalance) || !isFinite(minBalance) {
        return nil, AccountError{message: "inputs are invalid"}
    }
    
    if initialBalance < 0 || minBalance < 0 {
        return nil, NegativeAmountError{}
    }
    
	if initialBalance < minBalance {
	    return nil, InsufficientFundsError{minBalance: minBalance}
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
    if a == nil {
        return AccountError{message: "account cannot be nil"}
    }
    
    if !isFinite(amount) {
        return ExceedsLimitError{}
    }
    
    if amount < 0 {
        return NegativeAmountError{}
    }
    
    if amount > MaxTransactionAmount {
        return ExceedsLimitError{}
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
    if a == nil {
        return AccountError{message: "account cannot be nil"}
    }
    
    if !isFinite(amount) {
        return ExceedsLimitError{}
    }
    
    if amount < 0 {
        return NegativeAmountError{}
    }
    
    if amount > MaxTransactionAmount {
        return ExceedsLimitError{}
    }
    
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if a.Balance - amount < a.MinBalance {
        return InsufficientFundsError{minBalance: a.MinBalance}
    }
    
    a.Balance -= amount
    
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
    if a == nil || target == nil {
        return AccountError{message: "invalid source or target account"}
    }
    
    if withdrawalErr := a.Withdraw(amount); withdrawalErr != nil {
        return withdrawalErr
    }
    
    if depositErr := target.Deposit(amount); depositErr != nil {
        a.mu.Lock()
        a.Balance += amount
        a.mu.Unlock()
        return depositErr
    }
    
	return nil
}

func isFinite(value float64) bool {
    return !math.IsNaN(value) && !math.IsInf(value, 0)
}
