package challenge7

import (
	"fmt"
	"sync"
)

/*
  Challenge 7: Bank Account with Error Handling
	Problem Statement
	Implement a simple banking system with proper error handling. You'll create a BankAccount struct that manages balance operations and implements appropriate error handling.

	Requirements
	Implement a BankAccount struct that has the following fields:

	ID (string): Unique identifier for the account
	Owner (string): Name of the account owner
	Balance (float64): Current balance of the account
	MinBalance (float64): Minimum balance that must be maintained
	Implement the following methods:

	NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error): Constructor that validates input parameters
	Deposit(amount float64) error: Adds money to the account
	Withdraw(amount float64) error: Removes money from the account
	Transfer(amount float64, target *BankAccount) error: Transfers money from one account to another
	You must implement custom error types:

	InsufficientFundsError: When withdrawal/transfer would bring balance below minimum
	NegativeAmountError: When deposit/withdraw/transfer amount is negative
	ExceedsLimitError: When deposit/withdrawal amount exceeds your defined limits
	AccountError: A general bank account error with appropriate subtypes
	Function Signatures
	// Constructor
	func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error)

	// Methods
	func (a *BankAccount) Deposit(amount float64) error
	func (a *BankAccount) Withdraw(amount float64) error
	func (a *BankAccount) Transfer(amount float64, target *BankAccount) error

	// Error types
	type AccountError struct {
		// Implement custom error type with appropriate fields
	}

	type InsufficientFundsError struct {
		// Implement custom error type with appropriate fields
	}

	type NegativeAmountError struct {
		// Implement custom error type with appropriate fields
	}

	type ExceedsLimitError struct {
		// Implement custom error type with appropriate fields
	}

	// Each error type should implement the Error() string method
	func (e *AccountError) Error() string
	func (e *InsufficientFundsError) Error() string
	func (e *NegativeAmountError) Error() string
	func (e *ExceedsLimitError) Error() string
	Constraints
	All amounts must be valid values (non-negative).
	Withdrawals/transfers cannot bring account balance below the minimum balance.
	Define a reasonable limit for deposits and withdrawals (e.g., $10,000).
	Error messages should be descriptive and include relevant information.
	All operations should be thread-safe (use proper synchronization mechanisms).



*/

// Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

type BankAccount struct {
	mu         sync.Mutex
	ID         string
	Owner      string
	Balance    float64
	MinBalance float64
}

type ExceedsLimitError struct {
	problem string
}
type NegativeAmountError struct {
	problem string
}
type InsufficientFundsError struct {
	problem string
}
type AccountError struct {
	problem string
}

func (a *AccountError) Error() string {
	return fmt.Sprint(a.problem)
}
func (a *ExceedsLimitError) Error() string {
	return fmt.Sprint(a.problem)
}
func (a *NegativeAmountError) Error() string {
	return fmt.Sprint(a.problem)
}
func (a *InsufficientFundsError) Error() string {
	return fmt.Sprint(a.problem)
}

func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if len(id) == 0 {
		return nil, &AccountError{
			problem: "id can't have zero length",
		}
	}
	if len(owner) == 0 {
		return nil, &AccountError{
			problem: "owner can't have zero length",
		}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			problem: "initial balance cannot be negative",
		}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{
			problem: "minimum balance cannot be negative",
		}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			problem: "initial balance cannot be smaller than minBalance",
		}
	}
	return &BankAccount{ID: id, Owner: owner, Balance: initialBalance, MinBalance: minBalance}, nil
}

func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{
			problem: "amount to be deposited cannot be negative",
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			problem: fmt.Sprintf("deposit amount exceeds limit of %.2f", MaxTransactionAmount),
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.Balance += amount
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{
			problem: "amount to be withdrawn cannot be negative",
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			problem: fmt.Sprintf("withdrawal amount exceeds limit of %.2f", MaxTransactionAmount),
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			problem: fmt.Sprintf("withdrawal of %.2f would bring balance (%.2f) below minimum (%.2f)", amount, a.Balance, a.MinBalance),
		}
	}

	a.Balance -= amount
	return nil
}

func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if target == nil {
		return &AccountError{
			problem: "target bank account cannot be nil",
		}
	}
	if amount < 0 {
		return &NegativeAmountError{
			problem: "transfer amount cannot be negative",
		}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			problem: fmt.Sprintf("transfer amount exceeds limit of %.2f", MaxTransactionAmount),
		}
	}
	if a.ID == target.ID {
		return &AccountError{
			problem: "cannot transfer to the same account",
		}
	}

	// Deadlock prevention: Establish consistent lock ordering based on unique ID strings
	if a.ID < target.ID {
		a.mu.Lock()
		target.mu.Lock()
	} else {
		target.mu.Lock()
		a.mu.Lock()
	}
	defer a.mu.Unlock()
	defer target.mu.Unlock()

	// Safe to check balance under lock
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			problem: fmt.Sprintf("transfer of %.2f would bring balance (%.2f) below minimum (%.2f)", amount, a.Balance, a.MinBalance),
		}
	}

	a.Balance -= amount
	target.Balance += amount
	return nil
}

func main() {
	// Simple local test verifying thread safety and custom errors
	acc1, _ := NewBankAccount("acc_1", "Alice", 1000, 100)
	acc2, _ := NewBankAccount("acc_2", "Bob", 500, 50)

	err := acc1.Transfer(200, acc2)
	if err == nil {
		fmt.Printf("Transfer successful! Alice Balance: %.2f, Bob Balance: %.2f\n", acc1.Balance, acc2.Balance)
	}

	err = acc1.Withdraw(1000)
	if err != nil {
		fmt.Println("Expected Error caught:", err)
	}
}
