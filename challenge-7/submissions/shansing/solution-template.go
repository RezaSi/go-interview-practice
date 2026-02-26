// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
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

//var allAccounts []BankAccount

// Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	// Implement this error type
}

func (e *AccountError) Error() string {
	// Implement error message
	return "General Account Error!"
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	// Implement this error type
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return "Insufficient Funds!"
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	// Implement this error type
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	return "Negative Amount!"
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	// Implement this error type
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return "Exceeds Limit!"
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if len(id) == 0 {
		return nil, &AccountError{}
	}
	//if account, _ := getAccountById(id); account != nil {
	//	return nil, &AccountError{}
	//}
	if len(owner) == 0 {
		return nil, &AccountError{}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{}
	}
	balance := initialBalance
	if err := checkBalance(balance, minBalance); err != nil {
		return nil, err
	}
	BankAccount := &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    balance,
		MinBalance: minBalance,
	}
	return BankAccount, nil
}

//func getAccountById(id string) (*BankAccount, error) {
//	for account := range allAccounts {
//		if allAccounts[account].ID == id {
//			return &allAccounts[account], nil
//		}
//	}
//	return nil, &AccountError{}
//}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling// Implement deposit functionality with proper error handling
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	if amount < 0 {
		return &NegativeAmountError{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	newBalance := a.Balance + amount
	if err := checkBalance(newBalance, a.MinBalance); err != nil {
		return err
	}
	a.Balance = newBalance
	return nil
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Implement withdrawal functionality with proper error handling
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	if amount < 0 {
		return &NegativeAmountError{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	newBalance := a.Balance - amount
	if err := checkBalance(newBalance, a.MinBalance); err != nil {
		return err
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	a.Balance = newBalance
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	// Implement transfer functionality with proper error handling
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	if amount < 0 {
		return &NegativeAmountError{}
	}
	if a.ID < target.ID {
		a.mu.Lock()
		defer a.mu.Unlock()
		target.mu.Lock()
		defer target.mu.Unlock()
	} else if a.ID > target.ID {
		target.mu.Lock()
		defer target.mu.Unlock()
		a.mu.Lock()
		defer a.mu.Unlock()
	} else {
		return &AccountError{}
	}
	aBalance := a.Balance - amount
	if err := checkBalance(aBalance, a.MinBalance); err != nil {
		return err
	}
	targetBalance := target.Balance + amount
	if err := checkBalance(targetBalance, target.MinBalance); err != nil {
		return err
	}
	a.Balance = aBalance
	target.Balance = targetBalance
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	if amount < 0 {
		return &NegativeAmountError{}
	}
	return nil
}

func checkBalance(balance float64, minBalance float64) error {
	if balance < 0 {
		return &NegativeAmountError{}
	}
	if balance < minBalance {
		return &InsufficientFundsError{}
	}
	if balance > MaxTransactionAmount {
		return &ExceedsLimitError{}
	}
	return nil
}
