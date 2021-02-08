package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		counts[input.Text()]++
	}

	for i, n := range counts {
		if n > 1 {
			fmt.Println(fmt.Sprintf("word: '%s' entered %d times", i, n))
		}
	}

}
