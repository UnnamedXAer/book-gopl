package main

import (
	"fmt"
	"os"
)

func main() {
	// s := []rune{0x30d7, 0x30ed, 0x30b0, 0x30e9, 0x30e0}
	s := []rune{0x30d711ff, 0x30ed, 0x30b0, 0x30e9, 0x30e0}

	fmt.Fprintf(os.Stdout, "%s\n", string([]byte(string(s))[0]))
	fmt.Fprintf(os.Stdout, "%s\n", string([]byte(string(s))[1]))
	fmt.Fprintf(os.Stdout, "%s\n", string([]byte(string(s))[2]))
	fmt.Fprintf(os.Stdout, "%s\n", string(s))
	fmt.Fprintf(os.Stdout, "% x\n", string(s))
	fmt.Fprintf(os.Stdout, "% x\n", s)
}
