package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	double(142)

	bigSlowOperation()
}

func bigSlowOperation() {
	defer trace("bigSlowOperation")()
	// ... a lot of work
	time.Sleep(10 * time.Second) // simulate out "a lot of work"
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s\n", msg)
	return func() {
		log.Printf("exit %s (%s)\n", msg, time.Since(start))
	}
}

func double(x int) (result int) {
	defer func() {
		fmt.Printf("double(%d) = %d", x, result)
	}()
	return x + x
}
