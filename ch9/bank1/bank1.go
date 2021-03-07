// Package bank1 provides a concurrency-safe bank with one account.
package bank1

var (
	deposits = make(chan int) // send amount to deposit
	balances = make(chan int) // receive balance
)

// Deposit signals to update balance with given amount
func Deposit(amount int) { deposits <- amount }

// Balance returns current balances
func Balance() int { return <-balances }

func teller() {
	var balance int // balance is confined ti teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
