package bank3

import "testing"

func BenchmarkBalance(b *testing.B) {

	Deposit(123456)
	n := b.N
	sema := make(chan struct{}, n)

	x := func() {
		Balance()
		sema <- struct{}{}
	}

	go func() {
		for i := 0; i < n; i++ {
			go x()
		}
	}()

	for i := 0; i < b.N; i++ {
		<-sema
	}

}
