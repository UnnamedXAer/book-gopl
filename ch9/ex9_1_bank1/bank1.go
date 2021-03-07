// Package main exceed package bank1
package main

import "fmt"

var (
	deposits  = make(chan int) // send amount to deposit
	balances  = make(chan int) // receive balance
	withdraws = make(chan struct {
		amount int
		ch     chan<- response
	}) // receive balance
)

// Deposit signals to update balance with given amount
func Deposit(amount int) { deposits <- amount }

// Balance returns current balances
func Balance() int { return <-balances }

type response struct {
	ok      bool
	balance int
}

func Withdraw(amount int) (balance int, ok bool) {

	res := make(chan response)

	// go func() {
	withdraws <- struct {
		amount int
		ch     chan<- response
	}{
		amount: amount,
		ch:     res,
	}
	// }()
	fmt.Println("waiting for Withdraw result")

	select {
	case result := <-res:
		balance = result.balance
		ok = result.ok
	}
	fmt.Println("Withdraw result", ok, balance)

	return balance, ok
}

// teller is a monitor goroutine
func teller() {
	var balance int // balance is confined ti teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case withraw := <-withdraws:
			fmt.Println(1)
			if (withraw.amount < 0) || ((balance - withraw.amount) < 0) {
				fmt.Println(2)
				withraw.ch <- response{
					ok:      false,
					balance: balance,
				}
				continue
			}
			fmt.Println(3)

			balance -= withraw.amount
			withraw.ch <- response{
				ok:      true,
				balance: balance,
			}
			fmt.Println(4)
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
