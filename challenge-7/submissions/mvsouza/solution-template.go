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
	return e.message
}

type InsufficientFundsError struct {
	AttemptedAmount float64
	CurrentBalance  float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds to withdraw %.2f, current balance is %.2f", e.AttemptedAmount, e.CurrentBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	message string
}

func (e *NegativeAmountError) Error() string {
	return e.message
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return "When deposit/withdrawal amount exceeds your defined limits"
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" {
		return nil, &AccountError{"Account id can't be empty"}
	} else if owner == "" {
		return nil, &AccountError{"Account Owner can't be empty"}
	} else if initialBalance < 0 {
		return nil, &NegativeAmountError{"Account initial balance must be greater than 0"}
	} else if minBalance < 0 {
		return nil, &NegativeAmountError{"Account min balance must be greater than 0"}
	} else if minBalance > initialBalance {
		return nil, &InsufficientFundsError{AttemptedAmount: initialBalance, CurrentBalance: minBalance}
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
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount < 0 {
		return &NegativeAmountError{"Deposit amount must greater than 0"}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	a.Balance += amount
	return nil
}

// CheckBalanceAvailable checks if the account has enough balance for a withdrawal or transfer.
func (a *BankAccount) CheckBalanceAvailable(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{"It's only allowed to withdraw positive numbers"}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	} else if amount > a.Balance || a.MinBalance > a.Balance-amount {
		return &InsufficientFundsError{AttemptedAmount: amount, CurrentBalance: a.Balance}
	}
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.CheckBalanceAvailable(amount); err != nil {
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
		return &AccountError{"Can't transfer to the same account"}
	}
	if a.ID < target.ID {
		a.mu.Lock()
		target.mu.Lock()
	} else {
		target.mu.Lock()
		a.mu.Lock()
	}
	defer a.mu.Unlock()
	defer target.mu.Unlock()
	if err := a.CheckBalanceAvailable(amount); err != nil {
		return err
	}
	a.Balance -= amount
	target.Balance += amount
	return nil
}
