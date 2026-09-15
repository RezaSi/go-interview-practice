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

// MaxTransactionAmount Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	Message string
	Field   string
}

func (e *AccountError) Error() string {
	// Implement error message
	return fmt.Sprintf("AccountError: %s: %s", e.Field, e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Message string
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Insufficient funds: %s", e.Message)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Message string
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount: %s", e.Message)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Message string
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Exceeds limit of money: %s", e.Message)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, &AccountError{
			Message: "Empty ID",
			Field:   "ID",
		}
	}
	if owner == "" {
		return nil, &AccountError{
			Message: "Empty Owner",
			Field:   "Owner",
		}
	}
	if initialBalance <= 0 {
		return nil, &NegativeAmountError{
			Message: "Initial balance must be greater than zero",
		}
	}
	if minBalance <= 0 {
		return nil, &NegativeAmountError{
			Message: "Min balance must be greater than zero",
		}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			Message: "Initial balance must be greater than minimum balance",
		}
	}
	// Implement account creation with validation
	newBankAccount := BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}

	return &newBankAccount, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if amount >= MaxTransactionAmount {
		return &ExceedsLimitError{Message: "Transaction amount is too big"}
	}
	if amount < 0 {
		return &NegativeAmountError{Message: "Transaction amount is negative"}
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
	if amount >= MaxTransactionAmount {
		return &ExceedsLimitError{Message: "Transaction amount is too big"}
	}
	if amount < 0 {
		return &NegativeAmountError{Message: "Transaction amount is negative"}
	}
	if amount > a.Balance || a.Balance-amount < 0 || a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Message: "This operation led to insufficient funds"}
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
	wErr := a.Withdraw(amount)
	if wErr != nil {
		return wErr
	}

	dErr := target.Deposit(amount)
	if dErr != nil {
		return dErr
	}

	return nil
}
