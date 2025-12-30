// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"strings"
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
	message string

}

func (e *AccountError) Error() string {
	// Implement error message
	return fmt.Sprintf("Error with the account - %v", e.message )
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Balance    float64
	Amount  float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
    return fmt.Sprintf("insufficient funds: balance $%.2f, attempted to withdraw $%.2f", e.Balance, e.Amount)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
    return fmt.Sprintf("Amount $%.2f is less than 0", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Amount float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("Amount $%.2f exceeds MaxTransactionAmount $%.2f", e.Amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if len(strings.TrimSpace(id)) == 0{
	    return nil, &AccountError{"No ID given"}
	} 
	if len(strings.TrimSpace(owner)) == 0{
        return nil, &AccountError{"No owner given"}
	} 
	if  initialBalance < 0{
	    return nil, &NegativeAmountError{initialBalance}
	} 
	if  minBalance < 0{
	    return nil, &NegativeAmountError{minBalance}
	} 
	if  initialBalance < minBalance {
	    return nil, &InsufficientFundsError{Balance: initialBalance, Amount:minBalance  }
	} 
	return &BankAccount{		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount < 0{
	    return &NegativeAmountError{amount}
	}
	if amount > MaxTransactionAmount {
	     return &ExceedsLimitError{amount}
	}
    a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Implement withdrawal functionality with proper error handling
	if amount < 0 {
	    return &NegativeAmountError{amount}
	}
	if amount > MaxTransactionAmount {
	     return &ExceedsLimitError{amount}
	}
	if a.Balance - amount < a.MinBalance {
	    return &InsufficientFundsError{Balance:a.Balance,Amount:amount}
	}

	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	if amount < 0 {
	    return &NegativeAmountError{amount}
	}
	if amount > MaxTransactionAmount {
	     return &ExceedsLimitError{amount}
	}
	if a.Balance - amount < a.MinBalance {
	    return &InsufficientFundsError{Balance:a.Balance,Amount:amount}
	}

	a.Balance -= amount
	target.Balance += amount
	return nil
} 