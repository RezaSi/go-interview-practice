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

// ---------------- Custom Error Types ---------------- //

// AccountError is a general error type for bank account operations.
type AccountError struct {
	Operation string
	Reason    string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("Account error during %s: %s", e.Operation, e.Reason)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	AccountID string
	Attempted float64
	Balance   float64
	Min       float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Insufficient funds in account %s: tried %.2f, balance %.2f, min required %.2f",
		e.AccountID, e.Attempted, e.Balance, e.Min)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Operation string
	Amount    float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount for %s: %.2f", e.Operation, e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Operation string
	Amount    float64
	Limit     float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Transaction limit exceeded for %s: %.2f (limit %.2f)",
		e.Operation, e.Amount, e.Limit)
}

// ---------------- BankAccount Methods ---------------- //

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, &AccountError{"CreateAccount", "ID cannot be empty"}
	}
	if owner == "" {
		return nil, &AccountError{"CreateAccount", "Owner cannot be empty"}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{"InitialBalance", initialBalance}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{"MinBalance", minBalance}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{id, initialBalance, initialBalance, minBalance}
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
	if amount < 0 {
		return &NegativeAmountError{"Deposit", amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{"Deposit", amount, MaxTransactionAmount}
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
		return &NegativeAmountError{"Withdraw", amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{"Withdraw", amount, MaxTransactionAmount}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.ID, amount, a.Balance, a.MinBalance}
	}

	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{"Transfer", amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{"Transfer", amount, MaxTransactionAmount}
	}

	// Lock both accounts in consistent order to prevent deadlocks
	if a.ID < target.ID {
		a.mu.Lock()
		target.mu.Lock()
	} else {
		target.mu.Lock()
		a.mu.Lock()
	}
	defer a.mu.Unlock()
	defer target.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.ID, amount, a.Balance, a.MinBalance}
	}

	a.Balance -= amount
	target.Balance += amount
	return nil
}
