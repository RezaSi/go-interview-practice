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
	InvalidFieldKey string
	InvalidFieldValue string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf(
	    "Invalid value (%s) for field (%s)", 
	    e.InvalidFieldValue, 
	    e.InvalidFieldKey,
	)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// can adjust to requirements
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return "Funds insufficient"
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// can adjust to requirements
}

func (e *NegativeAmountError) Error() string {
	return "Amount is invalid (negetive)"
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// can adjust to requirements
}

func (e *ExceedsLimitError) Error() string {
	return "ExceedsLimitError"
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
    // ID or Owner invalid
	if id == "" {
	    return nil, &AccountError{
	        InvalidFieldKey: "ID", 
	        InvalidFieldValue: id,
	    }
	}
	if owner == "" {
	    return nil, &AccountError{
	        InvalidFieldKey: "Owner", 
	        InvalidFieldValue: owner,
	    }
	}
	
	if initialBalance < 0 || minBalance < 0 {
	    return nil, &NegativeAmountError{}
	}
	// Min Balance requirement not met
	if initialBalance < minBalance{
	    return nil, &InsufficientFundsError{}
	}
	newAccountPtr := &BankAccount{
	    ID: id,
	    Owner: owner,
	    Balance: initialBalance,
	    MinBalance: minBalance,
	    mu: sync.Mutex{},
	}
	return newAccountPtr, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
    // validate parameter alone
	if err := validateTxnAmount(amount); err != nil {
	    return err
	}
	// Lock before any account CRUD
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
    // validate parameter alone
	if err := validateTxnAmount(amount); err != nil {
	    return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance - amount < a.MinBalance {
	    return &InsufficientFundsError{}
	}
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// validate parameter alone
	if err := validateTxnAmount(amount); err != nil {
	    return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance - amount < a.MinBalance {
	    return &InsufficientFundsError{}
	}
	target.mu.Lock()
	defer target.mu.Unlock()
	a.Balance -= amount
	target.Balance += amount
	return nil
} 

// validateTxnAmount checks if an amount is within 0 to MaxTransactionAmount and returns an appropriate error otherwise
func validateTxnAmount(amount float64) error {
    if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{}
	} 
	if amount < 0 {
	    return &NegativeAmountError{}
	}
	return nil
}