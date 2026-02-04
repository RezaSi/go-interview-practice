// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
    "sync"
    "unsafe"
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

// AccountError is a general error type for bank account operations.
type AccountError struct {
	err string
}

func NewAccountError(err string) AccountError {
    return AccountError{
        err: err,
    }
}

func (e AccountError) Error() string {
	return "Account error: " + e.err
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {}

func (e InsufficientFundsError) Error() string {
	return "Insufficient funds error"
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {}

func (e NegativeAmountError) Error() string {
	return "Negative amount error"
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {}

func (e ExceedsLimitError) Error() string {
	return "Exceeds limit error"
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
    if len(id) == 0 {
        return nil, NewAccountError("empty id")
    }
    if len(owner) == 0 {
        return nil, NewAccountError("empty owner")
    }
    if initialBalance < 0 {
        return nil, NegativeAmountError{}
    }
    if minBalance < 0 {
        return nil, NegativeAmountError{}
    }
    if initialBalance < minBalance {
        return nil, InsufficientFundsError{}
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
    if amount < 0 {
        return NegativeAmountError{}
    }
    if amount > MaxTransactionAmount {
        return ExceedsLimitError{}
    }
    
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if a.Balance - amount < a.MinBalance {
        return InsufficientFundsError{}
    }
    
    a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
    if target == nil {
        return NewAccountError("target account is nil")
    }
    if a == target {
        return NewAccountError("cannot transfer to same account")
    }
    if amount < 0 {
        return NegativeAmountError{}
    }
    if amount > MaxTransactionAmount {
        return ExceedsLimitError{}
    }
    
	first, second := a, target
	if a.ID > target.ID ||
		(a.ID == target.ID && uintptr(unsafe.Pointer(a)) > uintptr(unsafe.Pointer(target))) {
		first, second = target, a
	}
    
    first.mu.Lock()
    defer first.mu.Unlock()
    second.mu.Lock()
    defer second.mu.Unlock()
    
    if a.Balance - amount < a.MinBalance {
        return InsufficientFundsError{}
    }
    
    a.Balance -= amount
    target.Balance += amount
    return nil
} 