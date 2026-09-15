package challenge7

import (
	"fmt"
	"sync"
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
	ID    string
	Owner string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("account error: %s: %s", e.ID, e.Owner)
}

type InsufficientFundsError struct {
	Balance float64
	Amount  float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: balance $%.2f, attempted to withdraw $%.2f", e.Balance, e.Amount)
}

type NegativeAmountError struct {
	Amount float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("negative amount: $%.2f", e.Amount)
}

type ExceedsLimitError struct {
	Amount float64
	Limit  float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("amount $%.2f exceeds limit $%.2f", e.Amount, e.Limit)
}

func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" || owner == "" {
		return nil, &AccountError{
			ID:    id,
			Owner: owner,
		}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{Amount: initialBalance}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{Amount: minBalance}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			Balance: initialBalance,
			Amount:  minBalance,
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
	if err := validateAmount(amount); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Balance: a.Balance,
			Amount:  amount,
		}
	}

	a.Balance -= amount
	return nil
}

func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if a == target {
		return nil
	}

	if err := validateAmount(amount); err != nil {
		return err
	}

	first, second := a, target
	if first.ID > second.ID {
		first, second = second, first
	}

	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{
			Balance: a.Balance,
			Amount:  amount,
		}
	}

	a.Balance -= amount
	target.Balance += amount

	return nil
}

func validateAmount(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{Amount: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Amount: amount, Limit: MaxTransactionAmount}
	}
	return nil
}
