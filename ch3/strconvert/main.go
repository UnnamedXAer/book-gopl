package main

import (
	"fmt"
)

func main() {
	// for _, v := range os.Args[1:] {
	// 	n, _ := strconv.Atoi(v)
	// 	fmt.Printf("number: %d\n", n)
	// 	fmt.Printf("%q %b\n", "%b", n)
	// 	fmt.Printf("%q %d\n", "%d", n)
	// 	fmt.Printf("%q %U\n", "%U", n)
	// 	fmt.Printf("%q %x\n", "%x", n)
	// }

	for i := 0; i < 8; i++ {
		fmt.Printf("%d - > % 4d,\t %08b\n", i, 1<<i, 1<<i)
	}

}
