package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("program requires a string argument to be suplied")
	}
	s := os.Args[1]
	fmt.Println()
	s = expand(s, foo)
	fmt.Println(s)
}

func expand(s string, f func(string) string) string {
	return strings.ReplaceAll(s, "$foo", f(".X."))
}

func foo(foo string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(100)
	// if n == 0 {
	// 	n++
	// 	n++
	// }
	s := make([]string, n, n)
	fmt.Println("n = ", n)
	return strings.Join(s, foo)
}
