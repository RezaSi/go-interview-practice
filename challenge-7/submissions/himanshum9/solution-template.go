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
	return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	minimum float64
	amount float64
	balance float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("the minimum balance should be remained which is &%0.2f available balance:&%0.2f, Transaction Amount:%0.2f",e.minimum,e.balance,e.amount)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	amount float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf("amount can not be negetive:&%.2f",e.amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	amount float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("amount limit exceeded from 10000, Amount given: &%.2f",e.amount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" || owner == "" {
	    return nil,&AccountError{Message:"Id or Owner is missing"}
	}
	
	if minBalance < 0 {
		return nil, &NegativeAmountError{amount: minBalance}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{amount: initialBalance}
	}
	
	if initialBalance < minBalance {
	    return nil, &InsufficientFundsError{minimum:minBalance,amount:0,balance:initialBalance}
	}
	    return &BankAccount{ID:id,Owner:owner,Balance:initialBalance,MinBalance:minBalance,mu:sync.Mutex{}}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0 {
	    return &NegativeAmountError{amount:amount}
	}
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{amount:amount}
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
        return &NegativeAmountError{amount:amount}
    }
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{amount:amount}
    }
    if a.Balance - amount < a.MinBalance {
        return &InsufficientFundsError{minimum:a.MinBalance,amount:amount,balance:a.Balance}
    }
    a.mu.Lock()
    a.Balance -= amount
    a.mu.Unlock()
	// Implement withdrawal functionality with proper error handling
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	if amount < 0 {
        return &NegativeAmountError{amount:amount}
    }
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{amount:amount}
    }
    if a.Balance - amount < a.MinBalance {
        return &InsufficientFundsError{minimum:a.MinBalance,amount:amount,balance:a.Balance}
    }
    a.mu.Lock()
    a.Balance -= amount
    target.Balance += amount
    a.mu.Unlock()
	return nil
} 