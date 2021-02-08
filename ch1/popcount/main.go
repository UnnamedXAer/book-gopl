package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// pc[i] is the population count of i
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCountLoop returns the population count (number of set bits) of x.
func PopCountLoop(x uint64) int {

	var p byte

	for i := 0; i < 8; i++ {
		p += pc[byte(x>>i*8)]
	}

	return int(p)
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>0*8)] +
		pc[byte(x>>1*8)] +
		pc[byte(x>>2*8)] +
		pc[byte(x>>3*8)] +
		pc[byte(x>>4*8)] +
		pc[byte(x>>5*8)] +
		pc[byte(x>>6*8)] +
		pc[byte(x>>7*8)])
}

func main() {
	fmt.Println()
	x, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	out1 := PopCount(uint64(x))
	out2 := PopCountLoop(uint64(x))

	fmt.Println(out1)
	fmt.Println(out2)
}
