// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
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
	ID string
}

func (e *AccountError) Error() string {
	if e.ID == "" {
		return fmt.Sprintf(`This account cannot be created`)
	}
	return fmt.Sprintf(`This account is invalid: %v`, e.ID)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	amount  float32
	balance float32
}

func (e *InsufficientFundsError) Error() string {

	return fmt.Sprintf(`Attemped to withdraw %v but only has a balance of %v`, e.amount, e.balance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	amount float32
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return fmt.Sprintf(`Cannot withdraw, deposit or set a negative amount: %v`, e.amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	amount float32
}

func (e *ExceedsLimitError) Error() string {

	return fmt.Sprintf(`The amount %v exceeds the limit of this account at: %v`, e.amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" || owner == "" {
		return nil, &AccountError{}
	} else if initialBalance <= 0 || minBalance <= 0 {
		return nil, &NegativeAmountError{}
	} else if minBalance > initialBalance {
		return nil, &InsufficientFundsError{}
	}

	newAcct := BankAccount{ID:id, Owner: owner, Balance: initialBalance,  MinBalance: minBalance}
	
	return &newAcct, nil
}
// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	switch {
	case amount < 0:
		return &NegativeAmountError{}

	case amount > MaxTransactionAmount:
		return &ExceedsLimitError{}
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

	switch {
	case amount < 0:
		return &NegativeAmountError{}
		
		case amount > MaxTransactionAmount:
		return &ExceedsLimitError{}

	case amount > (a.Balance - a.MinBalance):
		return &InsufficientFundsError{}

	}

	a.Balance -= amount
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	switch {
	case amount < 0:
		return &NegativeAmountError{}

	case amount > MaxTransactionAmount:
		return &ExceedsLimitError{}
		
	case amount > (a.Balance - a.MinBalance):
		return &InsufficientFundsError{}

	case amount > MaxTransactionAmount:
		return &ExceedsLimitError{}
	}

	target.Balance += amount
	a.Balance -= amount

	return nil
}