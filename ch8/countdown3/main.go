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
	tick := time.Tick(1 * time.Second)
	fmt.Println("Commencing countdown. Press return to abort.")

	for countdown := 10; countdown > 0; countdown-- {
		select {
		case <-tick:
			fmt.Printf("\r% 4d", countdown)
		case <-abort:
			fmt.Println("lauching aborted")
			return
		}
	}
	launch()

}

func launch() {
	fmt.Println()
	fmt.Println(strings.ToUpper("launched!"))
}
