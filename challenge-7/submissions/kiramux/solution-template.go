// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
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
	Code      string
	Message   string
	AccountID string
}

func (e *AccountError) Error() string {
	if e.AccountID != "" {
		return fmt.Sprintf("[%s] AccountID: %s, %s", e.Code, e.AccountID, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	Code       string
	Message    string
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("[%s] %s, your balance is less than the min balance: %.2f", e.Code, e.Message, e.MinBalance)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	Code    string
	Message string
	Amount  float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("[%s] %s, provided number: %.2f", e.Code, e.Message, e.Amount)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	Code    string
	Message string
	Amount  float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("[%s] %s, provided number: %.2f, the limit is %.2f", e.Code, e.Message, e.Amount, MaxTransactionAmount)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	// Determine the validity of the parameters.

	// Validate accountID
	if id == "" {
		return nil, &AccountError{
			Code:      "INVALID_ACCOUNT_ID",
			Message:   "account ID cannot be empty",
			AccountID: id,
		}
	}

	// Validate owner
	if owner == "" {
		return nil, &AccountError{
			Code:      "INVALID_OWNER",
			Message:   "owner name cannot be empty",
			AccountID: id,
		}
	}

	// Validate initial balance
	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			Code:    "INVALID_INITIAL_BALANCE",
			Message: "initial balance cannot be negative",
			Amount:  initialBalance,
		}
	}

	// Validate minimum balance
	if minBalance < 0 {
		return nil, &NegativeAmountError{
			Code:    "INVALID_MIN_BALANCE",
			Message: "min balance cannot be negative",
			Amount:  minBalance,
		}
	}

	// Compare initial balance and minimum balance
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			Code:       "INSUFFICIENT_FUND",
			Message:    fmt.Sprintf("the initialBalance: %.2f is less than minBalance: %.2f", initialBalance, minBalance),
			MinBalance: minBalance,
		}
	}

	return &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{
			Code:    "INVALID_DEPOSIT_AMOUNT",
			Message: "deposit amount cannot be negative",
			Amount:  amount,
		}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Code:    "EXCEED_LIMIT",
			Message: "deposit amount cannot exceed the limit",
			Amount:  amount,
		}
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
	if amount < 0 {
		return &NegativeAmountError{
			Code:    "INVALID_WITHDRAW_AMOUNT",
			Message: "withdraw amount cannot be negative",
			Amount:  amount,
		}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Code:    "EXCEED_LIMIT",
			Message: "withdraw amount cannot exceed the limit",
			Amount:  amount,
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	remain := a.Balance - amount
	if remain < a.MinBalance {
		return &InsufficientFundsError{
			Code:       "INSUFFICIENT_FUNDS",
			Message:    "account balance cannot be less than min amount",
			MinBalance: a.MinBalance,
		}
	}
	a.Balance = remain
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{
			Code:    "INVALID_TRANSFER_AMOUNT",
			Message: "transfer amount cannot be negative",
			Amount:  amount,
		}
	} else if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			Code:    "EXCEED_LIMIT",
			Message: "transfer amount cannot exceed the limit",
			Amount:  amount,
		}
	}

	// check target account is valid or not
	switch target {
	case nil:
		return &AccountError{
			Code:      "INVALID_TARGET_ACCOUNT",
			Message:   "target account is not existed",
			AccountID: "",
		}
	case a:
		return &AccountError{
			Code:      "INVALID_TARGET_ACCOUNT",
			Message:   "target account cannot be the from account",
			AccountID: a.ID,
		}
	}

	// The lock order is determined by the account ID number
	var first, second *BankAccount
	if a.ID < target.ID {
		first = a
		second = target
	} else if a.ID > target.ID {
		first = target
		second = a
	} else {
		// a.ID == target.ID but a != target (duplicate IDs)
		return &AccountError{
			Code:      "DUPLICATE_ACCOUNT_ID",
			Message:   "source and target accounts have duplicate IDs",
			AccountID: a.ID,
		}
	}

	first.mu.Lock()
	second.mu.Lock()
	defer second.mu.Unlock()
	defer first.mu.Unlock()

	remain := a.Balance - amount
	if remain < a.MinBalance {
		return &InsufficientFundsError{
			Code:       "INSUFFICIENT_FUNDS",
			Message:    "account balance cannot be less than min amount",
			MinBalance: a.MinBalance,
		}
	}
	a.Balance = remain
	target.Balance += amount
	return nil
}
