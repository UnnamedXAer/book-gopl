package main

import (
	"bufio"
	"fmt"
	"os"
)

type occurrence struct {
	cnt   int
	files []string
}

func main() {
	counts := make(map[string]occurrence)

	fmt.Println(os.Args[0])
	fmt.Println()
	args := os.Args[1:]
	filesCount := len(args)

	if filesCount == 0 {
		countLines(os.Stdin, counts, "")
	} else {
		cwd, _ := os.Getwd()
		for _, arg := range args {
			path := cwd + "/" + arg
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(fmt.Sprintf("error: %q, file: %q", err, path))
				continue
			}
			defer file.Close()
			countLines(file, counts, arg)
			fmt.Println("close err: ", err)
		}
	}

	for word, n := range counts {
		if n.cnt > 1 {
			fmt.Println(fmt.Sprintf("word: '%s' entered %d times in files: %v", word, n.cnt, n.files))
		}
	}
}

func countLines(file *os.File, counts map[string]occurrence, f string) {

	input := bufio.NewScanner(file)

	for input.Scan() {
		txt := input.Text()
		o := counts[txt]
		o.cnt++
		if f != "" {
			var exist bool
			for _, v := range o.files {
				if v == f {
					exist = true
					break
				}
			}
			if exist == false {

				o.files = append(o.files, f)
			}
		}
		counts[txt] = o
	}
}
