
package challenge7

import (
	"fmt"
	"strings"
	"sync"
)

type BankAccount struct {
	ID         string
	Owner      string
	Balance    float64
	MinBalance float64
	mu         sync.Mutex 
}

func (a *BankAccount) GetBalance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Balance
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
	return &BankAccount{ID: id, Owner: owner, Balance: initialBalance, MinBalance: minBalance}, nil
}

func (a *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{Value: amount}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: a.Balance + amount}
	}
	a.Balance += amount
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{Value: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: amount}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Balance: a.Balance, MinBalance: a.MinBalance}
	}
	a.Balance -= amount
	return nil
}

func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{Value: amount}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: amount}
	}
	if target == nil {
		return fmt.Errorf("target account cannot be nil")
	}
	if a == target {
		return fmt.Errorf("cannot transfer to the same account")
	}
	first, second := a, target
	if strings.Compare(a.ID, target.ID) > 0 {
		first, second = target, a
	}
	first.mu.Lock()
	defer first.mu.Unlock()

	second.mu.Lock()
	defer second.mu.Unlock()

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{Balance: a.Balance, MinBalance: a.MinBalance}
	}
	if target.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{Value: target.Balance + amount}
	}
	a.Balance -= amount
	target.Balance += amount
	return nil
}

func main() {
	account1, err := NewBankAccount("ACC001", "Alice", 1000.0, 100.0)
	if err != nil {
		fmt.Printf("Error creating account1: %v\n", err)
		return
	}
	account2, err := NewBankAccount("ACC002", "Bob", 500.0, 50.0)
	if err != nil {
		fmt.Printf("Error creating account2: %v\n", err)
		return
	}
	if err := account1.Deposit(200.0); err != nil {
		fmt.Printf("Error depositing: %v\n", err)
	}
	if err := account1.Withdraw(100.0); err != nil {
		fmt.Printf("Error withdrawing: %v\n", err)
	}
	if err := account1.Transfer(300.0, account2); err != nil {
		fmt.Printf("Error transferring: %v\n", err)
	}
	fmt.Printf("Account1 final balance: %.2f\n", account1.GetBalance())
	fmt.Printf("Account2 final balance: %.2f\n", account2.GetBalance())
}