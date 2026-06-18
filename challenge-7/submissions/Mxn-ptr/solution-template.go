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
	Message string
}

func (e *AccountError) Error() string {
	// Implement error message
	return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Amount     float64
	Balance    float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Insufficient funds: balance $%.2f with minimum balance $%.2f, attempted to withdraw $%.2f.\n",
		e.Balance, e.MinBalance, e.Amount)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount: $%.2f is negative\n", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Amount float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Exceeds limit: $%.2f is higher thant the limit %.2f\n", e.Amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, &AccountError{"Id must be non empty"}
	}
	if owner == "" {
		return nil, &AccountError{"Owner must be non empty"}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{initialBalance}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{minBalance}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			Amount:     0,
			Balance:    initialBalance,
			MinBalance: minBalance,
		}
	}
	bankAccount := BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}
	return &bankAccount, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{Amount: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Amount: amount}
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
		return &NegativeAmountError{Amount: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Amount: amount}
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	updatedBalance := a.Balance - amount
	if updatedBalance < a.MinBalance {
		return &InsufficientFundsError{
			Amount:     amount,
			Balance:    a.Balance,
			MinBalance: a.MinBalance,
		}
	}

	a.Balance = updatedBalance

	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if target == nil {
		return &AccountError{Message: "target account is nil"}
	}

	if err := a.Withdraw(amount); err != nil {
		return err
	}

	if err := target.Deposit(amount); err != nil {
		a.Deposit(amount)
		return err
	}

	return nil
}
