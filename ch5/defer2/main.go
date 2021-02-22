package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	defer printStack()
	err := f(3)

	fmt.Println("- ", err)
}

func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

func f(x int) (err error) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x = 0

	defer func() {
		fmt.Printf("deferred. %d\n", x)
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
		}
	}()

	return f(x - 1)
}
