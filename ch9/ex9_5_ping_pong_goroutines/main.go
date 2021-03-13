package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)
	ch2 := make(chan int)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			select {
			case x := <-ch:
				ch2 <- x + 1
			case <-done:
				wg.Done()
				return
			}

		}
	}()

	go func() {
		for {
			select {
			case x := <-ch2:
				ch <- x + 1
			case <-done:
				wg.Done()
				return
			}
		}

	}()

	ch <- 0
	t := time.Now()
	select {
	case <-time.After(1 * time.Second):
		close(done)
	}

	log.Println(time.Since(t))

	select {
	case x := <-ch:
		log.Println("x:", x)
	case x := <-ch2:
		log.Println("x2:", x)
	}
	wg.Wait()
}
