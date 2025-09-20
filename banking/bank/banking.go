package accounts

import (
	"errors"
	"fmt"
)

type Account struct {
	owner   string
	balance int
}

var errorNoMoney = errors.New("can't withdraw, insufficient funds")

func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

func (a *Account) Deposit(amount int) {
	a.balance += amount
}

func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errorNoMoney
	}
	a.balance -= amount
	return nil
}

func (a Account) Balance() int {
	return a.balance
}

func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

func (a Account) Owner() string {
	return a.owner

}
func (a Account) String() string {
	return fmt.Sprintf("%s's account. Balance: %d", a.owner, a.balance)
}
