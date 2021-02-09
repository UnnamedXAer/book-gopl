package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := strconv.ParseInt(os.Args[1], 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	// f = 0x12
	fmt.Fprintf(os.Stdout, "%v, %08b", f, f)
}
