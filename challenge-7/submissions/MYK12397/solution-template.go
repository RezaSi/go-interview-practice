// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	// Add any other necessary imports
	"fmt"
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
	ID        string
	Message   string
	Operation string
}

func (e AccountError) Error() string {
	// Implement error message
	msg := fmt.Sprintf("account error [%s]: %s - operation: %s",
		e.ID, e.Message, e.Operation)
	return msg
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
	ID              string
	RequestedAmount float64
	AvailableAmount float64
	Message         string
}

func (e InsufficientFundsError) Error() string {
	// Implement error message
	msg := fmt.Sprintf("insufficient funds in account %s: requested amount %.2f, available %.2f", e.ID, e.RequestedAmount, e.AvailableAmount)
	return msg
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
	Amount    float64
	Operation string
	Message   string
}

func (e NegativeAmountError) Error() string {
	// Implement error message
	msg := fmt.Sprintf("negative amount error: %.2f for operation '%s'", e.Amount, e.Operation)
	return msg
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
	ID        string
	Amount    float64
	Limit     float64
	LimitType string
}

func (e ExceedsLimitError) Error() string {
	// Implement error message
	msg := fmt.Sprintf("amount %.2f exceeds %s limit of %.2f for account id %s", e.Amount, e.LimitType, e.Limit, e.ID)
	return msg
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	negErr := NegativeAmountError{}
	negErrMsg := negErr.Error()
	negErr.Amount = initialBalance
	negErr.Message = negErrMsg
	if initialBalance < 0 || minBalance < 0 {
		return nil, negErr
	}

	insErr := InsufficientFundsError{}
	insErrMsg := insErr.Error()
	insErr.ID = id
	insErr.Message = insErrMsg
	if initialBalance > 0 && minBalance > 0 && initialBalance < minBalance {
		return nil, insErr
	}
	accErr := AccountError{}
	accErrMsg := accErr.Error()
	accErr.Message = accErrMsg
	if id == "" {
		return nil, accErr
	}
	if owner == "" {
		return nil, accErr
	}

	if id != "" && owner != "" && initialBalance >= 0.0 && minBalance >= 0.0 && initialBalance >= minBalance {
		return &BankAccount{
			ID:         id,
			Owner:      owner,
			Balance:    initialBalance,
			MinBalance: minBalance,
		}, nil
	}

	return nil, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	// Implement deposit functionality with proper error handling
	if amount > MaxTransactionAmount {
		limitErr := ExceedsLimitError{
			ID:     a.ID,
			Amount: a.Balance,
			Limit:  MaxTransactionAmount,
		}

		return &limitErr
	}
	if amount < 0 {
		negErr := NegativeAmountError{
			Amount:    amount,
			Operation: "deposit",
		}
		return &negErr
	}
	a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Implement withdrawal functionality with proper error handling
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount < 0 {
		NegErr := NegativeAmountError{}
		NegErr.Amount = amount
		NegErr.Operation = "withdraw"
		return NegErr
	}
	if amount > MaxTransactionAmount {
		limitErr := ExceedsLimitError{
			ID:     a.ID,
			Amount: a.Balance,
			Limit:  MaxTransactionAmount,
		}
		return limitErr
	}
	if (a.Balance - amount) < a.MinBalance {
		insErr := InsufficientFundsError{
			ID:              a.ID,
			RequestedAmount: amount,
			AvailableAmount: a.Balance,
		}
		return insErr
	}
	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount < 0 {
		NegErr := NegativeAmountError{}
		NegErr.Amount = amount
		NegErr.Operation = "withdraw"
		return NegErr
	}

	if amount > MaxTransactionAmount {
		limitErr := ExceedsLimitError{}
		return &limitErr
	}
	if (a.Balance - amount) < a.MinBalance {
		insErr := InsufficientFundsError{}
		return &insErr
	}
	a.Balance -= amount
	target.Balance += amount

	return nil
}
