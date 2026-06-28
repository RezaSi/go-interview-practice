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

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	MissingFields   []string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("Error, missing fields %v. Cannot create account.", e.MissingFields)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Amt     float64
	Min     float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Error, cannot transact $%v since that would lower account below minimum $%v.", e.Amt, e.Min)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amt     float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Error, cannot transact negative amount $%v.", e.Amt)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
    Amt     float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Error, transaction amount $%v is larger than maximum transaction amount $%v.", e.Amt, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	missingFields := []string{}
	if id == "" {
	    missingFields = append(missingFields, "id")
	} else if owner == ""{
	    missingFields = append(missingFields, "owner")
	}
	if len(missingFields) != 0 {
        return nil, &AccountError{ missingFields }
	}
    
    if initialBalance < 0.0 {
        return nil, &NegativeAmountError{ initialBalance }
    }
    if minBalance < 0.0 {
        return nil, &NegativeAmountError{ minBalance }
    }
    if initialBalance < minBalance {
        return nil, &InsufficientFundsError{ 
            Amt: initialBalance,
            Min: minBalance,
        }
    }
    
    acct := BankAccount{
        ID: id,
        Owner: owner,
        Balance: initialBalance,
        MinBalance: minBalance,
    }
	return &acct, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if err := isAmtValid(amount); err != nil {
	    return err
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
	if err := isAmtValid(amount); err != nil {
	    return err
	}
	
	a.mu.Lock()
    defer a.mu.Unlock()
	
	if err := doesWithdrawalGoBelowMinBalance(amount, a.Balance, a.MinBalance); err != nil {
	    return err
	}
    a.Balance -= amount
    
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
    if err := isAmtValid(amount); err != nil {
	    return err
	}
	
    a.mu.Lock()
    defer a.mu.Unlock()
    target.mu.Lock()
    defer target.mu.Unlock()
    
	if err := doesWithdrawalGoBelowMinBalance(amount, a.Balance, a.MinBalance); err != nil {
	    return err
	}
    a.Balance -= amount
    target.Balance += amount

    return nil
} 

func isAmtValid(amount float64) error {
    if amount < 0.0 {
	    return &NegativeAmountError{ amount }
	}
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{ amount }
	}
	return nil
}

func doesWithdrawalGoBelowMinBalance(withdrawalAmt, balance, minBalance float64) error {
    if expectedBal := balance - withdrawalAmt; expectedBal < minBalance {
	    return &InsufficientFundsError{ 
	        Amt: withdrawalAmt,
	        Min: minBalance,
        }
	}
	return nil
}