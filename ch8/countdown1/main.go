package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	fmt.Println("Commencing countdown.")

	tick := time.Tick(1 * time.Second)

	for countdown := 10; countdown > 0; countdown-- {
		fmt.Printf("\r% 4d", countdown)
		<-tick
	}

	launch()
}

func launch() {
	fmt.Println()
	fmt.Println(strings.ToUpper("launched!"))
}

// func main() {
// 	abort := make(chan struct{})

// 	go func() {
// 		os.Stdin.Read(make([]byte, 1)) //read a single byte
// 		abort <- struct{}{}
// 	}()
// }
