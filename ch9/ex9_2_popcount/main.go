package main

import (
	"fmt"
	"sync"
)

var mu = sync.Once{}
var pc [256]byte

func PopCount(x uint64) int {
	mu.Do(func() {
		fmt.Println("im doing it")
		for i := range pc {
			pc[i] = pc[i/2] + byte(i&1)
		}
	})

	return int(
		pc[byte(x>>0*8)] +
			pc[byte(x>>1*8)] +
			pc[byte(x>>2*8)] +
			pc[byte(x>>3*8)] +
			pc[byte(x>>4*8)] +
			pc[byte(x>>5*8)] +
			pc[byte(x>>6*8)] +
			pc[byte(x>>7*8)])
}

func main() {

	fmt.Println(PopCount(34))
	fmt.Println(PopCount(2 << 8))
	fmt.Println(uint64(2 << 8))

}
