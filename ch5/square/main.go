package main

import "fmt"

func main() {
	f := squares()
	fmt.Println("f", f())
	fmt.Println("f", f())
	fmt.Println("f", f())
	f2 := squares()
	fmt.Println("f", f())
	fmt.Println("f2", f2())
	fmt.Println("f2", f2())
	fmt.Println("f", f())
	fmt.Println("f2", f2())
}

func squares() func() int {
	var x int

	return func() int {
		x++
		return x * x
	}
}
