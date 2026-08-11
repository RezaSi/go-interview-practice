// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"fmt"
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
	// Implement this error type
	Message string
}

func (e *AccountError) Error() string {
	// Implement error message
	return fmt.Sprintf("AccountError: %s", e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	Balance float64
	Amount float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("InsufficientFundsError: Current Balance is %.1f, but attempt to withdraw/transfer %.1f", e.Balance, e.Amount)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	Namount float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf("NegativeAmountError: Current balance is Invalid %.1f", e.Namount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	Amount float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("ExceedsLimitError: Withdraw/Deposit of %.1f, exceeds limit of %.1f", e.Amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" || owner == "" {
	    return nil, &AccountError{
	        Message: "Empty ID/Owner",
	    }
	}
	if initialBalance > MaxTransactionAmount {
	    return nil, &ExceedsLimitError{
	        Amount: initialBalance, 
	    }
	}
	
	if initialBalance < 0 {
	    return nil, &NegativeAmountError{
	        Namount: initialBalance,
	    }
	}
	if minBalance < 0 {
	    return nil, &NegativeAmountError{
	        Namount: minBalance,
	    }
	}

	if initialBalance < minBalance {
	    return nil, &InsufficientFundsError{
	        Balance: initialBalance,
	        Amount: minBalance,
	    }
	}
	
	NewBankAccount := &BankAccount{
	    ID: id,
	    Owner: owner,
	    Balance: initialBalance,
	    MinBalance: minBalance,
	}
	
	return NewBankAccount, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0 {
	    return &NegativeAmountError{
	        Namount: amount,
	    }
	}
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Amount: amount,
	    } 
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance = a.Balance + amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Implement withdrawal functionality with proper error handling
	if amount < 0 {
	    return &NegativeAmountError{
	        Namount: amount,
	    }
	}
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Amount: amount,
	    } 
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if (a.Balance - amount) < a.MinBalance {
	    return &InsufficientFundsError{
	        Balance: a.Balance,
	        Amount: amount,
	    }
	}
	a.Balance = a.Balance - amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	if amount < 0 {
	    return &NegativeAmountError{
	        Namount: amount,
	    }
	}
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Amount: amount,
	    } 
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if (a.Balance - amount) < a.MinBalance {
	    return &InsufficientFundsError{
	        Balance: a.Balance,
	        Amount: amount,
	    }
	}
	target.mu.Lock()
	defer target.mu.Unlock()
	
	target.Balance = target.Balance + amount
	a.Balance = a.Balance - amount
	return nil
} 