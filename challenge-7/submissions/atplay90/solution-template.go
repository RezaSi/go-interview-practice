// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"errors"
	"fmt"
	"strings"
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
	Field string
	Message string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("account error on %s: %s", e.Field, e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Message string
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds error: %s", e.Message)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Message string
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("negative amount error: %s", e.Message)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Message string
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("exceeds limit error: %s", e.Message)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	
	if strings.TrimSpace(id) == "" {
		return nil, &AccountError{
			Field:   "id",
			Message: "account ID cannot be empty",
		}
	}
	
	if strings.TrimSpace(owner) == "" {
		return nil, &AccountError{
			Field:   "owner",
			Message: "owner name cannot be empty",
		}
	}
	
	if initialBalance < 0 {
		return nil, &NegativeAmountError{
		    Message: "initial balance cannot be negative",
		}
	}
	
	if minBalance < 0 {
		return nil, &NegativeAmountError{
		    Message: "minimum balance cannot be negative",
		}
	}
	
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
		    Message: fmt.Sprintf("initial balance %.2f cannot be less than minimum balance %.2f", initialBalance, minBalance),
		}
	}
	
	return &BankAccount{
		ID:             strings.TrimSpace(id),
		Owner:          strings.TrimSpace(owner),
		Balance:        initialBalance,
		MinBalance: minBalance,
	}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if amount < 0.0 {
        return &NegativeAmountError{
            Message: "amount to deposit cannot be negative",
        }
    }
    
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{
            Message: "amount to deposit exceeds max transaction amount",
        }
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
    
    if amount < 0.0 {
        return &NegativeAmountError{
            Message: "amount to withdraw cannot be negative",
        }
    }
    
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{
            Message: "amount to withdraw exceeds max transaction amount",
        }
    }
    
    if a.Balance - amount < a.MinBalance {
        return &InsufficientFundsError{
            Message: "amount to withdraw brings the balance below the minimum required balance",
        }
    }
    
    a.Balance -= amount

	return nil

}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
    if amount < 0.0 {
        return &NegativeAmountError{
            Message: "amount to transfer cannot be negative",
        }
    }
    
    if amount > MaxTransactionAmount {
        return &ExceedsLimitError{
            Message: "amount to transfer exceeds max transaction amount",
        }
    }
    
    if a.ID == target.ID {
		return errors.New("cannot transfer money to the same account")
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
    
    if a.Balance - amount < a.MinBalance {
        return &InsufficientFundsError{
            Message: "amount to transfer brings the balance below the minimum required balance",
        }
    }
    
    a.Balance -= amount
	target.Balance += amount

	return nil
} 