package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]

	fmt.Println("program: ", os.Args[0])

	t := time.Now()
	fmt.Println(strings.Join(args, " "))
	duration := time.Since(t)
	fmt.Println("duration of the join: ", duration)

	t = time.Now()
	s, sep := "", ""
	for i := 0; i < len(args); i++ {
		s = sep + s
		sep = " "
	}
	fmt.Println(s)
	duration = time.Since(t)
	fmt.Println("duration of the loop: ", duration)

	t = time.Now()
	s1, sep1 := "", " "
	count := len(args)
	if count > 0 {

		s1 = args[0]
	}
	for i := 1; i < count; i++ {
		s1 = sep1 + s1
	}
	fmt.Println(s)
	duration = time.Since(t)
	fmt.Println("duration of the loop optimized: ", duration)

	for i, v := range args {
		fmt.Println(fmt.Sprintf("arg[%d] is: %s", i, v))
	}
}
