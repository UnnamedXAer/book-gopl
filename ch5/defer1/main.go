package main

import "fmt"

func main() {
	f(4)
	f(0)
}

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x = 0

	defer func() {
		fmt.Printf("deferred. %d\n", x)
	}()
	f(x - 1)
}
