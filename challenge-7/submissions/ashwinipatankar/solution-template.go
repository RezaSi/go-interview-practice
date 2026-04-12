package challenge7 


import (
	"fmt"
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

// Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	Code    int
	Message string
}

func (e *AccountError) Error() string {
	// Implement error message
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Balance float64
	Amount  float64
}

func (e *InsufficientFundsError) Error() string {
	// Implement error message
	return fmt.Sprintf("insufficient funds: balance $%.2f, attempted to withdraw $%.2f",
		e.Balance, e.Amount)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Amount  float64
	Balance *float64
}

func (e *NegativeAmountError) Error() string {
	// Implement error message
	msg := fmt.Sprintf("negative amount: amount $%.2f", e.Amount)

	return msg
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Limit  float64
	Amount float64
}

func (e *ExceedsLimitError) Error() string {
	// Implement error message
	return fmt.Sprintf("exceeds limit: limit: $%.2f, trying to deposit-withdraw $%.2f",
		e.Limit, e.Limit-e.Amount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Implement account creation with validation
	if id == "" || owner == "" {
		return nil, &AccountError{1, "invalid account parameters"}
	}

	if initialBalance < 0.0 || minBalance < 0 {
		return nil, &NegativeAmountError{initialBalance, &minBalance}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{initialBalance, minBalance}
	}
	
	return &BankAccount{ID: id, Owner: owner, Balance: initialBalance, MinBalance: minBalance}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	// Implement deposit functionality with proper error handling
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{MaxTransactionAmount, amount}
	}

	if amount < 0 {
		return &NegativeAmountError{amount, &a.Balance}
	}

	if a.Balance+amount < a.MinBalance {
		return &InsufficientFundsError{a.MinBalance, amount}
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
	// Implement withdrawal functionality with proper error handling
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{MaxTransactionAmount, amount}
	}

	if amount < 0 {
		return &NegativeAmountError{amount, &a.Balance}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.Balance, amount}
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
	// Implement transfer functionality with proper error handling
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{MaxTransactionAmount, amount}
	}

	if amount < 0 {
		return &NegativeAmountError{amount, &a.Balance}
	}

	if a.Balance < amount {
		return &InsufficientFundsError{a.Balance, amount}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{a.Balance, amount}
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance -= amount
	target.Balance += amount

	return nil
}
