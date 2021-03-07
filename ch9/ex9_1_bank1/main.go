package main

import (
	"fmt"
)

func main() {

	for i := 0; i < 1000; i++ {
		go Deposit(i)
	}

	wantedAmount := 424933
	balance, ok := Withdraw(wantedAmount)
	if ok {
		fmt.Println("you successfully drawn your cash, current balace is:", balance)
	} else {
		fmt.Println("no luck today, you couldn't draw your cash, current balace is:", balance)
		fmt.Println("you would end up with balance:", balance-wantedAmount)
	}
}
