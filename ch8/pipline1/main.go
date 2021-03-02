package main

import (
	"log"
	"time"
)

func main() {

	c := make(chan int)
	cSqr := make(chan int)
	go counter(c)
	go square(c, cSqr)

	print(cSqr)
}

func counter(c chan<- int) {
	log.Println("counter", "- start")
	for i := 0; i < 10; i++ {
		log.Println("counter - emit: ", i)
		c <- i
		time.Sleep(100 * time.Millisecond)
	}
	close(c)
	log.Println("counter - end")

}

func square(c <-chan int, cSqr chan<- int) {
	log.Println("square ", "- start")
	for n := range c {
		sqr := n * n
		time.Sleep(30 * (time.Millisecond))
		log.Println("square - emit: ", sqr)
		cSqr <- sqr
	}
	close(cSqr)
	log.Println("square - end")

}

func print(cSqr chan int) {
	log.Println("print ", "- start")
	for {
		n, ok := <-cSqr
		if ok == false {
			break
		}
		log.Printf("print Sqr: %d", n)
	}
	log.Println("print - end")

}
