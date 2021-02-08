package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unnamedxaer/book-gopl/ch1/tempconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stdout, "\ncf: %v", err)
			continue
		}

		// fmt.Printf("\n\t %s\t %s\t %s", c, f, k)
		f := tempconv.Fahrenheit(t)
		c := tempconv.Celsius(t)
		fmt.Printf("%s = %s,\t %s = %s\n", f, tempconv.FToC(f), c, tempconv.CToF(c))

	}
}
