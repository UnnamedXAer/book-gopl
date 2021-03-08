package bank3

import (
	"sync"
	"time"
)

var (
	muRW   sync.RWMutex
	balace float64
)

func Deposit(amount float64) {
	muRW.Lock()
	deposit(amount)
	muRW.Unlock()
}

func Balance() float64 {
	muRW.RLock()
	time.Sleep(10 * time.Millisecond)
	b := balace
	muRW.RUnlock()
	return b
}

func Withdraw(amount float64) bool {
	muRW.Lock()
	defer muRW.Unlock()
	deposit(-amount)
	if balace < 0 {
		deposit(amount)
		return false
	}
	return true
}

// this function requires that the lock be held.
func deposit(amount float64) {
	balace += amount
}
