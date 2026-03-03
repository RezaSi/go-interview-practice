// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"fmt"
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
	ID    string
	Owner string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("invalid account: ID %s, Owner %s", e.ID, e.Owner)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Balance    float64
	MinBalance float64
	Amount     float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: balance %.2f becomes less then min balance %.2f",
		e.Balance-e.Amount, e.MinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("negative amount %.2f", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Amount float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("exceeds limit: amount %.2f, limit %.2f",
		e.Amount, MaxTransactionAmount)
}

type SameAccountTransferError struct {
	ID string
}

func (e *SameAccountTransferError) Error() string {
	return fmt.Sprintf("error while transferring to the same account with ID %s", e.ID)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	account := &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}
	err := account.validate()
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	err := validateAmountLimit(amount)
	if err != nil {
		return err
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
	err := validateAmountLimit(amount)
	if err != nil {
		return err
	}
	err = a.validateMinBalance(amount)
	if err != nil {
		return err
	}
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if a.ID == target.ID {
		return &SameAccountTransferError{a.ID}
	}
	if a.ID < target.ID {
		a.mu.Lock()
		defer a.mu.Unlock()
		target.mu.Lock()
		defer target.mu.Unlock()
	} else {
		target.mu.Lock()
		defer target.mu.Unlock()
		a.mu.Lock()
		defer a.mu.Unlock()
	}
	err := validateAmountLimit(amount)
	if err != nil {
		return err
	}
	err = a.validateMinBalance(amount)
	if err != nil {
		return err
	}
	a.Balance -= amount
	target.Balance += amount
	return nil
}

func (a *BankAccount) validate() error {
	if a.ID == "" || a.Owner == "" {
		return &AccountError{a.ID, a.Owner}
	}
	if a.Balance < 0 {
		return &NegativeAmountError{a.Balance}
	}
	if a.MinBalance < 0 {
		return &NegativeAmountError{a.MinBalance}
	}
	if a.Balance < a.MinBalance {
		return &InsufficientFundsError{a.Balance, a.MinBalance, 0}
	}
	return nil
}

func (a *BankAccount) validateMinBalance(amount float64) error {
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.Balance, a.MinBalance, amount}
	}
	return nil
}

func validateAmountLimit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{amount}
	}
	return nil
}
