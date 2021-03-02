package main

import (
	"log"
)

func main() {

	c := make(chan int, 10)
	cSqr := make(chan int, 5)
	go counter(c)
	go square(c, cSqr)

	print(cSqr)
}

func counter(c chan<- int) {
	log.Println("-C", "- start")
	for i := 0; i < 10; i++ {
		log.Println("-C - emit: ", i)
		c <- i
	}
	close(c)
	log.Println("-C - end")

}

func square(c <-chan int, cSqr chan<- int) {
	log.Println("- S ", "- start")
	for n := range c {
		sqr := n * n
		log.Println("- S - emit: ", sqr)
		cSqr <- sqr
	}
	close(cSqr)
	log.Println("- S - end")

}

func print(cSqr chan int) {
	log.Println("-   P ", "- start")
	for n := range cSqr {
		log.Printf("-   P - Sqr: %d", n)
	}
	log.Println("-   P - end")

}
