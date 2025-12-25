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

const MaxTransactionAmount = 10000.0

type AccountError struct {
	ID    string
	Owner string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("Cannot create Account with ID: %s, Owner: %s", e.ID, e.Owner)
}

type InsufficientFundsError struct {
	Balance    float64
	MinBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Insufficient balance: %f with min balance: %f", e.Balance, e.MinBalance)
}

type NegativeAmountError struct {
	Value float64
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative balance: %f", e.Value)
}

type ExceedsLimitError struct {
	Value float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Exceeds balance: %f", e.Value)
}

func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if len(owner) == 0 || len(id) == 0 {
		return nil, &AccountError{ID: id, Owner: owner}
	}
	if initialBalance < 0 {
		return nil, &NegativeAmountError{Value: initialBalance}
	}
	if minBalance < 0 {
		return nil, &NegativeAmountError{Value: minBalance}
	}
	if initialBalance > MaxTransactionAmount {
		return nil, &ExceedsLimitError{Value: initialBalance}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{Balance: initialBalance, MinBalance: minBalance}
	}
	var mu sync.Mutex
	return &BankAccount{ID: id, Owner: owner, Balance: initialBalance, MinBalance: minBalance, mu: mu}, nil
}

func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{Value: amount}
	}
	if a.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: a.Balance + amount}
	}
	a.mu.Lock()
	a.Balance += amount
	a.mu.Unlock()
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{Value: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: a.Balance - amount}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Balance: a.Balance, MinBalance: a.MinBalance}

	}
	a.mu.Lock()
	a.Balance -= amount
	a.mu.Unlock()
	return nil
}

func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{Value: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: a.Balance - amount}
	}
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Balance: a.Balance, MinBalance: a.MinBalance}
	}
	if target.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: a.Balance + amount}
	}
	a.mu.Lock()
	a.Balance -= amount
	target.Balance += amount
	a.mu.Unlock()
	return nil
}

func main() {
	account1, err := NewBankAccount("ACC001", "Alice", 1000.0, 100.0)
	if err != nil {
		fmt.Println("Aa")
	}
	account2, err := NewBankAccount("ACC002", "Bob", 500.0, 50.0)
	if err != nil {
		fmt.Println("Aa")
	}
	if err := account1.Deposit(200.0); err != nil {
		fmt.Println("Aa")
	}
	if err := account1.Withdraw(100.0); err != nil {
		fmt.Println("Aa")
	}
	if err := account1.Transfer(300.0, account2); err != nil {
		fmt.Println("Aa")
	}
	fmt.Println(account1.Balance)
}
