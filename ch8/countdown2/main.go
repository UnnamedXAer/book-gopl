package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	abort := make(chan struct{})

	go func() {
		os.Stdin.Read(make([]byte, 1)) //read a single byte
		abort <- struct{}{}
	}()
	// tick := time.Tick(1 * time.Second)
	fmt.Println("Commencing countdown. Press return to abort.")

	select {
	case <-time.After(10 * time.Second):
		// do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()

	// for countdown := 10; countdown > 0; countdown-- {
	// 	fmt.Printf("\r% 4d", countdown)
	// 	<-tick
	// }
	// launch()

}

func launch() {
	fmt.Println()
	fmt.Println(strings.ToUpper("launched!"))
}
