package main

import (
	"fmt"

	accounts "github.com/namlulu/banking/bank"
)

func main() {
	account := accounts.NewAccount("Nam")
	account.Deposit(500)

	err := account.Withdraw(2000)
	if err != nil {
		println(err.Error())
	}

	fmt.Println(account)
}
