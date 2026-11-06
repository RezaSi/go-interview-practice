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
	message string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("An error occured for account ID %q", e.message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	bankAccountId         string
	bankAccountBalance    float64
	bankAccountMinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Insufficient funds for account ID %q: new balance %f is below minimum balance %f", e.bankAccountId, e.bankAccountBalance, e.bankAccountMinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	bankAccountId string
	amount        float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount value for account ID %q: %f", e.bankAccountId, e.amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	bankAccountId string
	amount        float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Exceeding transaction limit for account ID %q: amount %f is above limit %f", e.bankAccountId, e.amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if len(id) == 0 {
		return nil, &AccountError{"Could not create bank account: empty ID"}
	}
	if len(owner) == 0 {
		return nil, &AccountError{"Could not create bank account: empty owner"}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{"", initialBalance}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{"", minBalance}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{"", initialBalance, minBalance}
	}
	return &BankAccount{id, owner, initialBalance, minBalance, sync.Mutex{}}, nil
}

func (a *BankAccount) verifyAmount(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{a.ID, amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{a.ID, amount}
	}
	return nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	err := a.verifyAmount(amount)
	if err != nil {
		return err
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
	err := a.verifyAmount(amount)
	if err != nil {
		return err
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.ID, a.Balance, a.MinBalance}
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
	err := a.verifyAmount(amount)
	if err != nil {
		return err
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.ID, a.Balance, a.MinBalance}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	target.mu.Lock()
	defer target.mu.Unlock()
	a.Balance -= amount
	target.Balance += amount
	return nil
}
