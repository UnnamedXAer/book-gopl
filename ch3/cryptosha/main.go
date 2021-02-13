package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

var shaType = flag.Int("n", 256, "sha256.Sum<n>; supported values are: 256, 512, 384")

func main() {

	flag.Parse()
	switch int(*shaType) {
	case 512:
		fmt.Printf("%x\n", sha512.Sum512([]byte("x")))
	case 384:
		fmt.Printf("%x\n", sha512.Sum384([]byte("x")))
	case 256:
		fmt.Printf("%x\n", sha256.Sum256([]byte("x")))
	default:
		fmt.Fprint(os.Stderr, fmt.Errorf("not supported value of n = %d", shaType))
	}
	// c1 := sha256.Sum256([]byte("x"))
	// c2 := sha256.Sum256([]byte("X"))
	// fmt.Printf("%x\n", c1)
	// fmt.Printf("%x\n", c2)
	// // var c3 [32]byte
	// var n int
	// for i := 0; i < 32; i++ {

	// 	n += bits.OnesCount(uint(c1[i] ^ c2[i]))
	// 	// fmt.Printf("%08b\n", c1[i])
	// 	// fmt.Printf("%08b\n", c2[i])
	// 	// fmt.Printf("%08b = %d\n", c1[i]^c2[i], bits.OnesCount(uint(c1[i]^c2[i])))
	// 	// fmt.Print("\n")
	// }
	// fmt.Println(n)
}
