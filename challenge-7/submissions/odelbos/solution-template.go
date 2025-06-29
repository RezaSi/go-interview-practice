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
    ID  string
    Op  string
    Msg string
}

func (e *AccountError) Error() string {
    return fmt.Sprintf("error, account: %s, op: %s, msg: %s", e.ID, e.Op, e.Msg)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
    ID      string
    Op      string
    Amount  float64
    Msg     string
}

func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf("error, account: %s, op: %s, amount: %f, msg: %s", e.ID, e.Op, e.Amount, e.Msg)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
    ID     string
    Op     string
    Amount float64
    Msg    string
}

func (e *NegativeAmountError) Error() string {
    return fmt.Sprintf("error, account: %s, op: %s, amount: %f, msg: %s", e.ID, e.Op, e.Amount, e.Msg)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
    ID     string
    Op     string
    Amount float64
    Msg    string
}

func (e *ExceedsLimitError) Error() string {
    return fmt.Sprintf("error, account: %s, op: %s, amount: %f, msg: %s", e.ID, e.Op, e.Amount, e.Msg)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
    if id == "" {
        return nil, &AccountError{id, "create", "cannot create account without valid ID"}
    }
    if owner == "" {
        return nil, &AccountError{id, "create", "cannot create account without valid owner"}
    }

    if initialBalance < 0 {
        return nil, &NegativeAmountError{id, "create", initialBalance, "initial balance is negative"}
    }
    if minBalance < 0 {
        return nil, &NegativeAmountError{id, "create", minBalance, "minimum balance is negative"}
    }

    if initialBalance < minBalance {
       return nil, &InsufficientFundsError{id, "create", initialBalance, "balance < minimum balance"}
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
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{a.ID, "deposit", amount, fmt.Sprintf("exceed the limit of: %f", MaxTransactionAmount)}
    }
    if amount < 0 {
        return &NegativeAmountError{a.ID, "deposit", amount, "amount cannot be negative"}
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
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{a.ID, "deposit", amount, fmt.Sprintf("exceed the limit of: %f", MaxTransactionAmount)}
    }
    if amount < 0 {
        return &NegativeAmountError{a.ID, "deposit", amount, "amount cannot be negative"}
    }

    a.mu.Lock()
    if (a.Balance - amount < a.MinBalance) {
        return &InsufficientFundsError{a.ID, "create", amount, "balance - amount < minimum balance"}
    }
    a.Balance -= amount
    a.mu.Unlock()
    return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
    err := a.Withdraw(amount)
    if err != nil {
        return err
    }
    target.mu.Lock()
    target.Balance += amount
    target.mu.Unlock()
    return nil
} 