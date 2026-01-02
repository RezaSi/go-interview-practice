// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package challenge7

import (
	"sync"
	"fmt"
)

type BankAccount struct {
	ID         string
	Owner      string
	Balance    float64
	MinBalance float64
	mu         sync.Mutex
}

const (
	MaxTransactionAmount = 10000.0
)

type AccountError struct {
	AccountID string
	Operation string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("account %s: error during %s", e.AccountID, e.Operation)
}

type InsufficientFundsError struct {
	AccountError
	Balance    float64
	Amount     float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf(
		"account %s: insufficient funds (balance: $%.2f, needed: $%.2f, min: $%.2f)",
		e.AccountID, e.Balance, e.Amount, e.MinBalance,
	)
}

type NegativeAmountError struct {
	AccountError
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf(
		"account %s: negative amount: $%.2f",
		e.AccountID, e.Amount,
	)
}

type ExceedsLimitError struct {
	AccountError
	Amount float64
	Limit  float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf(
		"account %s: amount $%.2f exceeds limit $%.2f",
		e.AccountID, e.Amount, e.Limit,
	)
}

func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, &AccountError{
			AccountID: "unknown",
			Operation: "create",
		}
	}

	if owner == "" {
		return nil, &AccountError{
			AccountID: id,
			Operation: "create",
		}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{
			AccountError: AccountError{
				AccountID: id,
				Operation: "create",
			},
			Amount: initialBalance,
		}
	}

	if minBalance < 0 {
		return nil, &NegativeAmountError{
			AccountError: AccountError{
				AccountID: id,
				Operation: "create",
			},
			Amount: minBalance,
		}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			AccountError: AccountError{
				AccountID: id,
				Operation: "create",
			},
			Balance:    initialBalance,
			Amount:     minBalance,
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

func (a *BankAccount) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if amount < 0 {
		return &NegativeAmountError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "deposit",
			},
			Amount: amount,
		}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "deposit",
			},
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	}

	a.Balance += amount
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if amount < 0 {
		return &NegativeAmountError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "withdraw",
			},
			Amount: amount,
		}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "withdraw",
			},
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "withdraw",
			},
			Balance:    a.Balance,
			Amount:     amount,
			MinBalance: a.MinBalance,
		}
	}

	a.Balance -= amount
	return nil
}

func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if a.ID < target.ID {
		a.mu.Lock()
		defer a.mu.Unlock()
		target.mu.Lock()
		defer target.mu.Unlock()
	} else {
		target.mu.Lock()
		defer target.mu.Unlock()
		a.mu.Lock()
		defer a.mu.Unlock()
	}

	if amount < 0 {
		return &NegativeAmountError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "transfer",
			},
			Amount: amount,
		}
	}

	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "transfer",
			},
			Amount: amount,
			Limit:  MaxTransactionAmount,
		}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			AccountError: AccountError{
				AccountID: a.ID,
				Operation: "transfer",
			},
			Balance:    a.Balance,
			Amount:     amount,
			MinBalance: a.MinBalance,
		}
	}

	a.Balance -= amount
	target.Balance += amount

	return nil
}