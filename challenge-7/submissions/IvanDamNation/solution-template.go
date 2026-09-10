// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
    "fmt"
    "math"
	"sync"
)

// BankAccount represents a bank account with balance management and minimum balance requirements.
type BankAccount struct {
	ID         string
	Owner      string
	Balance    float64  // Better use int64 or decimal
	MinBalance float64  // Can cause calc inaccuracy or NaN errors (IEEE 754)
	mu         sync.Mutex // For thread safety
}

// Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	ID string
	Err error
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("account ID %s error: %v", e.ID, e.Err)
}

func (e *AccountError) Unwrap() error {
    return e.Err
}

// ValidationError occurs when input is invalid
type ValidationError struct {
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
    Amount float64
	Balance float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("not enough funds: amount (%.2f) will bring balance (%.2f) below minimum (%.2f)", e.Amount, e.Balance, e.MinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("invalid input: negative number (%.2f)", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
    Amount float64
    Limit float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("limit exceeded: transaction amount (%.2f) is bigger than max limit (%.2f)", e.Amount, e.Limit)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initBalance, minBalance float64) (*BankAccount, error) {
    
    if id == "" {
        return nil, &AccountError{
            ID: id,
            Err: &ValidationError{
                Message: "ID cannot be empty",
            },
        }
    }
    
    if owner == "" {
        return nil, &AccountError{
            ID: id,
            Err: &ValidationError{
                Message: "owner cannot be empty",
            },
        }
    }
    
    if math.IsNaN(initBalance) || math.IsNaN(minBalance) {
        return nil, &AccountError {
            ID: id,
            Err: &ValidationError{
                Message: "initial balance or minimum balance cannot be NaN",
            },
        }
    }
    
    if math.IsInf(initBalance, 0) || math.IsInf(minBalance, 0) {
        return nil, &AccountError{
            ID: id,
            Err: &ValidationError{
                Message: "initial balance or minimum balance cannot be Inf",
            },
        }
    }
    
    if initBalance < 0 {
        return nil, &NegativeAmountError{
            Amount: initBalance,
        }
    }
    
    if minBalance < 0 {
        return nil, &NegativeAmountError{
            Amount: minBalance,
        }
    }
    
    if initBalance < minBalance {
        return nil, &InsufficientFundsError{
            Amount: 0,
            Balance: initBalance,
            MinBalance: minBalance,
        }
    }
    
	newAcc := &BankAccount{
	    ID: id,
	    Owner: owner,
	    Balance: initBalance,
	    MinBalance: minBalance,
	}
	
	return newAcc, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
	if err := a.deposit(amount); err != nil {
	    return err
	}
	
	return nil
}

func (a *BankAccount) deposit(amount float64) error {
    if math.IsNaN(amount) {
        return &AccountError{
            ID: a.ID,
            Err: &ValidationError{
                Message: "transaction amount cannot be NaN",
            },
        }
    }
    
    if math.IsInf(amount, 0) {
        return &AccountError{
            ID: a.ID,
            Err: &ValidationError{
                Message: "transaction amount cannot be Inf",
            },
        }
    }
    
    if amount < 0 {
	    return &NegativeAmountError{
	        Amount: amount,
	    }
	}
	
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Amount: amount,
	        Limit: MaxTransactionAmount,
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
    
	if err := a.withdraw(amount); err != nil {
	    return err
	}
	
	return nil
}

func (a *BankAccount) withdraw(amount float64) error {
    if math.IsNaN(amount) {
        return &AccountError{
            ID: a.ID,
            Err: &ValidationError{
                Message: "transaction amount cannot be NaN",
            },
        }
    }
    
    if math.IsInf(amount, 0) {
        return &AccountError{
            ID: a.ID,
            Err: &ValidationError{
                Message: "transaction amount cannot be Inf",
            },
        }
    }
    
    if amount < 0 {
	    return &NegativeAmountError{
	        Amount: amount,
	    }
	}
	
	if amount > MaxTransactionAmount {
	    return &ExceedsLimitError{
	        Amount: amount,
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
	
	a.Balance -= amount
	
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
    if target == nil {
        return &AccountError{
            ID: a.ID,
            Err: &ValidationError{Message: "target account cannot be nil"},
        }
    }
    
    if a.ID == target.ID {
        return &AccountError{
            ID: a.ID,
            Err: &ValidationError{
                Message: "sender and target accounts are identical",
            },
        }
    }
    
    if a.ID < target.ID {
        a.mu.Lock()
	    target.mu.Lock()
    } else {
        target.mu.Lock()
        a.mu.Lock()
    }
	
	defer a.mu.Unlock()
	defer target.mu.Unlock()
	
	var err error
	if err = a.withdraw(amount); err != nil {
	    return err
	}
	
	if err = target.deposit(amount); err != nil {
	    a.deposit(amount)
	    return err
	}
	
	return nil
}
