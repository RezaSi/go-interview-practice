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
	message string
}

func (e *AccountError) Error() string {
	return e.message
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	id      string
	owner   string
	amtName string
	amount  float64
	min     float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Acc Id: %s, owner: %s - %s value $%.2f is below minimum $%.2f", e.id, e.owner, e.amtName, e.amount, e.min)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	id      string
	owner   string
	amtName string
	amount  float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Acc Id: %s, owner: %s - %s value $%.2f is negative", e.id, e.owner, e.amtName, e.amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	id      string
	owner   string
	amtName string
	amount  float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Acc Id: %s, owner: %s - %s value $%.2f exceeds max transaction limit $%.2f",
		e.id, e.owner, e.amtName, e.amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if len(strings.TrimSpace(id)) == 0 {
		return nil, &AccountError{"invalid id"}
	}
	if len(strings.TrimSpace(owner)) == 0 {
		return nil, &AccountError{"invalid owner"}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			id:      id,
			owner:   owner,
			amtName: "new account initial balance",
			amount:  initialBalance,
		}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{
			id:      id,
			owner:   owner,
			amtName: "new account minimum balance",
			amount:  minBalance,
		}
	} // "initial balance less than minimum"
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			id:      id,
			owner:   owner,
			amtName: "new account initial balance",
			amount:  initialBalance,
			min:     minBalance,
		}
	}

	acc := BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}
	return &acc, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// lock mutex
	a.mu.Lock()
	defer a.mu.Unlock()

	// validate amount
	if amount < 0 {
		return &NegativeAmountError{
			id:      a.ID,
			owner:   a.Owner,
			amtName: "deposit amount",
			amount:  amount,
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			id:      a.ID,
			owner:   a.Owner,
			amtName: "deposit amount",
			amount:  amount,
		}
	}

	// increment balance
	a.Balance += amount

	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// lock mutex
	a.mu.Lock()
	defer a.mu.Unlock()

	// validate amount
	if amount < 0 {
		return &NegativeAmountError{
			id:      a.ID,
			owner:   a.Owner,
			amtName: "withdrawal amount",
			amount:  amount,
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			id:      a.ID,
			owner:   a.Owner,
			amtName: "withdrawal amount",
			amount:  amount,
		}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			id:      a.ID,
			owner:   a.Owner,
			amtName: "remaining balance",
			amount:  a.Balance - amount,
			min:     a.MinBalance,
		}
	}

	// decrement balance
	a.Balance -= amount

	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// capture initial balances
	initSrcBal := a.Balance
	initTgtBal := target.Balance

	if err := a.Withdraw(amount); err != nil {
		// revert initial balances
		a.Balance = initSrcBal
		target.Balance = initTgtBal
		return err
	}
	if err := target.Deposit(amount); err != nil {
		// revert initial balances
		a.Balance = initSrcBal
		target.Balance = initTgtBal
		return err
	}
	return nil
}