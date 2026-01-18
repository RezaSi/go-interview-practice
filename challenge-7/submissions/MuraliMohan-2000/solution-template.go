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
	AccountId string
	Message   string
}

func (e *AccountError) Error() string {
	if e.AccountId != "" {
		return fmt.Sprintf("account :%s : %s", e.AccountId, e.Message)
	}
	return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	AccountID  string
	Balance    float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf(
		"account %s: insufficient funds (balance=%.2f,minBalance=%.2f)",
		e.AccountID,
		e.Balance,
		e.MinBalance,
	)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("amount must be positive, got %.2f", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Amount float64
	Limit  float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf(
		"amount %.2f exceeds the allowed limit %.2f",
		e.Amount,
		e.Limit,
	)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if minBalance < 0 {
		return nil, &NegativeAmountError{
			Amount: minBalance,
		}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			Amount: initialBalance,
		}
	}

	if id == "" {
		return nil, &AccountError{
			Message: "account id cannot be empty",
		}
	}

	if owner == "" {
		return nil, &AccountError{
			Message: "owner name cannot be empty",
		}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			AccountID:  id,
			Balance:    initialBalance,
			MinBalance: minBalance,
		}
	}

	account := &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}

	return account, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0 {
		return &NegativeAmountError{
			Amount: amount,
		}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
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
	if amount < 0 {
		return &NegativeAmountError{
			Amount: amount,
		}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	} else if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			AccountID:  a.ID,
			Balance:    a.Balance,
			MinBalance: a.MinBalance,
		}

	}
	a.mu.Lock()
	a.Balance -= amount
	a.mu.Unlock()

	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{
			Amount: amount,
		}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	} else if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			AccountID:  a.ID,
			Balance:    a.Balance,
			MinBalance: a.MinBalance,
		}

	}
	a.mu.Lock()
	a.Balance -= amount
	a.mu.Unlock()

	target.mu.Lock()
	target.Balance += amount
	target.mu.Unlock()

	return nil
}
