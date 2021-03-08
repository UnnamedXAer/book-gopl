package bank2

var (
	sema    = make(chan struct{}, 1) // a binary semaphore guarding balance
	balance int
)

func Deposit(amaout int) {
	sema <- struct{}{} // cquire token
	balance += amaout
	<-sema // realease token
}

func Balance() int {
	sema <- struct{}{} // aquire token
	b := balance
	<-sema // release token
	return b
}
