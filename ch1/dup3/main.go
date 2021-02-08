package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	fmt.Println(os.Args[0])
	fmt.Println()
	args := os.Args[1:]
	filesCount := len(args)

	if filesCount == 0 {
		countLines(os.Stdin, counts)
	} else {
		cwd, _ := os.Getwd()
		for _, arg := range args {
			path := cwd + "/" + arg
			file, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(fmt.Sprintf("error: %q, file: %q", err, path))
				continue
			}

			lines := strings.Split(string(file), "\n")
			for _, line := range lines {
				counts[line]++
			}
		}
	}

	for i, n := range counts {
		if n > 1 {
			fmt.Println(fmt.Sprintf("word: '%s' entered %d times", i, n))
		}
	}
}

func countLines(file *os.File, counts map[string]int) {

	input := bufio.NewScanner(file)

	for input.Scan() {
		counts[input.Text()]++
	}
}
