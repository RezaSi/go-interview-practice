// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
    "fmt"
	"sync"
	"errors"
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
    AccountID string
    Operation string
    Reason error
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("%s failed for account %s: %v", e.Operation, e.AccountID, e.Reason)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Amount float64
	Balance float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("failed to change balance: operation would leave %.2f, while min balance is %.2f",
	    e.Balance - e.Amount, e.MinBalance,
	)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("amount must be a postitive number (recieved %.2f)", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
    Value float64
	Limit float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("transaction declined: amount %.2f exceeds limit %.2f", e.Value, e.Limit)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" || owner == "" {
        return nil, &AccountError{
            Operation: "create",
            Reason:    errors.New("id and owner cannot be empty"),
        }
    }
    
    if initialBalance < 0 {
        return nil, &NegativeAmountError{
            Amount: initialBalance,
        }
    }
    
    if minBalance < 0 {
        return nil, &NegativeAmountError{
            Amount: minBalance,
        }
    }
    
    if initialBalance < minBalance {
        return nil, &InsufficientFundsError{
            Amount: 0, Balance: initialBalance, MinBalance: minBalance,
        }
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
	    return &NegativeAmountError{
	        Amount: amount,
	    }
	}
	
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Value: amount,
	        Limit: MaxTransactionAmount,
	    }
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
	    return &NegativeAmountError{
	        Amount: amount,
	    }
	}
	
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Value: amount,
	        Limit: MaxTransactionAmount,
	    }
	}

	if a.Balance - amount < a.MinBalance {
	    return &InsufficientFundsError{
	        Amount: amount,
	        Balance: a.Balance,
	        MinBalance: a.MinBalance,
	    }
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
	err := a.Withdraw(amount)
	if err != nil {
	    return err
	}
	
	err = target.Deposit(amount)
	if err != nil {
	    return err
	}
	
	return nil
} 