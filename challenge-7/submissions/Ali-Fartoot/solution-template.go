// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"fmt"
	"strings"
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
	Message string
	Code    string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("account error [%s]: %s", e.Code, e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	RequestedAmount float64
	CurrentBalance  float64
	MinBalance      float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: requested %.2f, current balance %.2f, minimum balance %.2f",
		e.RequestedAmount, e.CurrentBalance, e.MinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("negative amount not allowed: %.2f", e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Amount float64
	Limit  float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("amount %.2f exceeds transaction limit of %.2f", e.Amount, e.Limit)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Validate input parameters
	if strings.TrimSpace(id) == "" {
		return nil, &AccountError{
			Message: "account ID cannot be empty",
			Code:    "INVALID_ID",
		}
	}

	if strings.TrimSpace(owner) == "" {
		return nil, &AccountError{
			Message: "account owner cannot be empty",
			Code:    "INVALID_OWNER",
		}
	}

	if minBalance < 0 {
		return nil, &NegativeAmountError{Amount: minBalance}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{Amount: initialBalance}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			RequestedAmount: 0,
			CurrentBalance:  initialBalance,
			MinBalance:      minBalance,
		}
	}

	return &BankAccount{
		ID:         strings.TrimSpace(id),
		Owner:      strings.TrimSpace(owner),
		Balance:    initialBalance,
		MinBalance: minBalance,
	}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Validate amount
	if amount < 0 {
		return &NegativeAmountError{Amount: amount}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	}

	// Thread-safe operation
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Validate amount
	if amount < 0 {
		return &NegativeAmountError{Amount: amount}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	}

	// Thread-safe operation
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if withdrawal would bring balance below minimum
	if amount > 0 && a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			RequestedAmount: amount,
			CurrentBalance:  a.Balance,
			MinBalance:      a.MinBalance,
		}
	}

	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Validate target account
	if target == nil {
		return &AccountError{
			Message: "target account cannot be nil",
			Code:    "INVALID_TARGET",
		}
	}

	// Validate amount
	if amount < 0 {
		return &NegativeAmountError{Amount: amount}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	}

	// Prevent self-transfer
	if a == target {
		return &AccountError{
			Message: "cannot transfer to the same account",
			Code:    "SELF_TRANSFER",
		}
	}

	// Lock both accounts in a consistent order to prevent deadlocks
	// Always lock the account with the smaller ID first
	var first, second *BankAccount
	if strings.Compare(a.ID, target.ID) < 0 {
		first, second = a, target
	} else {
		first, second = target, a
	}

	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()

	// Check if withdrawal would bring source balance below minimum
	if amount > 0 && a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			RequestedAmount: amount,
			CurrentBalance:  a.Balance,
			MinBalance:      a.MinBalance,
		}
	}

	// Perform the transfer
	a.Balance -= amount
	target.Balance += amount

	return nil
}

// GetBalance returns the current balance (thread-safe read)
func (a *BankAccount) GetBalance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Balance
}


// String returns a string representation of the account
func (a *BankAccount) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return fmt.Sprintf("Account[ID: %s, Owner: %s, Balance: %.2f, MinBalance: %.2f]",
		a.ID, a.Owner, a.Balance, a.MinBalance)
}