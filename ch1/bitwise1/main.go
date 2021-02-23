package main

import (
	"fmt"
	"math"
)

func main() {
	// var x uint8 = 1 << 1 //| 1<<5
	// var x uint8 = 1 << 5
	// var y uint8 = 1<<1 | 1<<2

	// fmt.Printf("%08b\n", x)
	// fmt.Printf("%08b\n", y)

	bitwise()
	// numbers()

}

func bitwise() {

	for i := 0; i < 10; i++ {
		fmt.Printf("%02d; %08b i\n", i, i)
		fmt.Printf("%02d; %08b (i-1)\n", (i - 1), (i - 1))
		// fmt.Printf("%02d; %08b i<<1\n", i<<1, i<<1)
		// fmt.Printf("%02d; %08b i|1\n", i|1, i|1)
		// fmt.Printf("%02d; %08b i&1\n", i&1, i&1)
		fmt.Printf("%02d; %08b i&^(i-1)\n\n", i&^(i-1), i&^(i-1))

	}
}

func numbers() {
	for x := 0; x < 8; x++ {
		xe := math.Exp(float64(x))
		fmt.Printf("x = %d e^x =%8.3f\n", x, xe)
	}

	fmt.Print(math.Sqrt(-1))

}
