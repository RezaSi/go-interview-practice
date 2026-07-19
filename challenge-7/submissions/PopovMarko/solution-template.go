// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"fmt"
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
	MaxTransactionAmount = 10_000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	Message string
}

func (e *AccountError) Error() string {
	return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Balance    float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {

	return fmt.Sprintf("Insufficient funds, balance %.2f, min balance %.2f", e.Balance, e.MinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount, not accepted: %.2f", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Amount   float64
	Restrict float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Max limit %.2f operation of %.2f not allowed", e.Restrict, e.Amount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Validate parameters
	// Check id not blank
	if id == "" {
		return nil, &AccountError{"ID is empty"}
	}

	// Check owner is not blank
	if owner == "" {
		return nil, &AccountError{"Owner is empty"}
	}

	// Check minBalance not negative
	if minBalance < 0 {
		return nil, &NegativeAmountError{minBalance}
	}

	// Check initialBalance not negative
	if initialBalance < 0 {
		return nil, &NegativeAmountError{initialBalance}
	}

	// Check initialBalance more than minBalance
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{initialBalance, minBalance}
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
	// Check amount not negative
	if amount < 0 {
		return &NegativeAmountError{amount}
	}

	// Check amount not exceeded operation limits
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{amount, MaxTransactionAmount}
	}

	// Update balance under Mutex
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Check amount not negative
	if amount < 0 {
		return &NegativeAmountError{amount}
	}

	// Check amount not exceeded its limit
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{amount, MaxTransactionAmount}
	}

	// Check for enough sum on balance under Mutex
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.Balance, a.MinBalance}
	}

	//Update balance under Mutex
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Check amount not negative
	if amount < 0 {
		return &NegativeAmountError{amount}
	}

	// Check amount not exceeded its limit
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{amount, MaxTransactionAmount}
	}

	//Guard against nil target
	if target == nil {
		return &AccountError{"Target account is nil"}
	}

	// Guard against self-transfer
	if a == target {
		return &AccountError{"Can't transfer to the same account"}
	}

	// Lock in canonical order to prevent deadlock
	first, second := a, target
	if uintptr(unsafe.Pointer(first)) > uintptr(unsafe.Pointer(second)) {
		first, second = second, first
	}

	// Lock Mutex for current and target account
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()

	// Check for enough sum for transfer
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.Balance, a.MinBalance}
	}

	// Update balance
	a.Balance -= amount
	target.Balance += amount
	return nil
}
