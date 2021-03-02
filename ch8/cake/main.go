package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

var sellDuration = time.Duration(time.Now().Unix()) / 20

func main() {
	ccap := 2
	fmt.Println(sellDuration, 30*time.Millisecond, sellDuration < 30*time.Millisecond)
	cBake := make(chan [3]string, ccap)
	cIced := make(chan [3]string, ccap)
	cReady := make(chan [3]string, ccap)

	fmt.Println("cap: ", cap(cBake), cap(cIced), cap(cReady))

	go bake(cBake)
	go icing(cBake, cIced)
	go icing(cBake, cIced)
	go icing(cBake, cIced)
	go icing(cBake, cIced)
	go icing(cBake, cIced)
	go icing(cBake, cIced)
	go icing(cBake, cIced)
	go inscribe(cIced, cReady)
	sell(cReady)
}

func bake(out chan<- [3]string) {
	for i := 0; i < 7; i++ {
		time.Sleep(30 * time.Millisecond)
		cake := [3]string{"baked #" + strconv.Itoa(i+1)}
		out <- cake
		log.Println(cake)
	}
	log.Println("closing cBake")
	close(out)
}

func icing(in <-chan [3]string, out chan<- [3]string) {
	for {
		time.Sleep(10000 * time.Millisecond)
		x, ok := <-in
		if ok == false {
			log.Println("closing cIced")
			close(out)
			break
		}

		inced := [3]string{x[0], "inced"}
		out <- inced
		log.Println("inced", inced[0])
	}

}

func inscribe(in <-chan [3]string, out chan<- [3]string) {
	for x := range in {
		time.Sleep(50 * time.Millisecond)
		ready := [3]string{x[0], x[1], "inscribed"}
		out <- ready
		log.Println("inscribed", ready[0])
	}
	log.Println("closing cReady")
	close(out)
}

func sell(in <-chan [3]string) {
	for cake := range in {
		time.Sleep(sellDuration)
		log.Println("sold", cake)
	}
	log.Println("Hurayy, all cakes sold!")
}
