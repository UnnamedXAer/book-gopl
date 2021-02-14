package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	// charCountWithCateg()
	wordCountFreq()
}

func charCount() {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error

		if err == io.EOF {
			break
		}

		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}

		counts[r]++
		utflen[n]++
	}

	fmt.Printf("rune\tcount\t\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\t\n", c, n)
	}

	fmt.Printf("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}

	fmt.Printf("\n%d invalid characters", invalid)
}

func charCountWithCateg() {
	counts := make(map[string]map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	// f, err := ioutil.ReadFile("t.txt")
	f, _ := os.Open("t.txt")
	defer f.Close()
	in := bufio.NewReader(f)
	// in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		var categ string
		if err == io.EOF {
			break
		}

		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}

		switch {

		case unicode.IsLetter(r):
			categ = "letter"
		case unicode.IsMark(r):
			categ = "mark"
		case unicode.IsNumber(r):
			categ = "number"
		case unicode.IsDigit(r):
			categ = "digit"
		case unicode.IsPunct(r):
			categ = "punct"
		case unicode.IsSpace(r):
			categ = "space"
		case unicode.IsSymbol(r):
			categ = "symbol"
		case unicode.IsControl(r):
			categ = "control"
		case unicode.IsGraphic(r):
			categ = "graphic"
		default:
			categ = "no-categ"
		}

		if counts[categ] == nil {
			counts[categ] = make(map[rune]int)
		}

		counts[categ][r]++
		utflen[n]++
	}

	for k, categ := range counts {

		fmt.Printf("\ncateg: %q\nrune\tcount\t\n", k)
		for c, n := range categ {
			fmt.Printf("%q\t%d\t\n", c, n)
		}
	}

	fmt.Printf("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}

	fmt.Printf("\n%d invalid characters", invalid)
}

func wordCountFreq() {
	counts := make(map[string]int)
	f, _ := os.Open("t.txt")
	defer f.Close()
	// b, _ := ioutil.ReadAll(f)
	in := bufio.NewScanner(f)
	in.Split(bufio.ScanWords)
	// bufio.ScanWords(b)
	for in.Scan() {
		s := in.Text() // returns rune, nbytes, error

		counts[s]++
	}

	fmt.Print("\nrune\tcount\t\n")
	for c, n := range counts {
		fmt.Printf("% 30q\t%d\t\n", c, n)
	}

}
