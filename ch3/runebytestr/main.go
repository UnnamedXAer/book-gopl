package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	// for _, v := range os.Args[1:] {
	// 	// if strings.Contains(v, ".") == false {
	// 	// 	fmt.Printf("rune: %q -> %q\n", v, commaByRune(v))
	// 	// 	fmt.Printf("byte: %q -> %q\n", v, commaByByte(v))
	// 	// }
	// 	// fmt.Printf("recursive: %q -> %q\n", v, comma(v))
	// 	// fmt.Printf("recursive: %q -> %q\n", v, commaFull(v))
	// }
	fmt.Printf("%q - %q -> %t\n", os.Args[1], os.Args[2], isAnagram(os.Args[1], os.Args[2]))
}

func commaByRune(s string) string {

	r := []rune(s)
	count := len(r)
	if count <= 3 {
		return s
	}

	sep := rune(',')
	scomma := []rune{r[0]}

	for i := 1; i < count; i++ {
		if (count-i)%3 == 0 {
			scomma = append(scomma, sep)
		}
		scomma = append(scomma, r[i])
	}

	// for i, v := range s {
	// 	if i > 0 && i%3 == 0 {
	// 		scomma = append(scomma, sep)
	// 	}
	// 	scomma = append(scomma, v)
	// }
	return string(scomma)
}

func commaByByte(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	var buf bytes.Buffer
	buf.WriteByte(s[0])
	for i := 1; i < n; i++ {
		if (n-i)%3 == 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(s[i])
	}
	return buf.String()
}

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {

	n := len(s)
	if n <= 3 {
		return s
	}

	return comma(s[:n-3]) + "," + s[n-3:]
}

func commaFull(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		fmt.Println(fmt.Errorf("too many dots %d", len(parts)))
	}
	n := len(parts[0])
	if n <= 3 {
		return strings.Join(parts, ".")
	}

	sign := s[0]
	var startIdx int
	if sign == '-' || sign == '+' {
		if n == 4 {
			return s
		}
		startIdx = 1
	}
	out := comma(parts[0][startIdx:])
	if len(parts) == 2 {
		out = out + "." + parts[1]
	}
	if startIdx == 1 {
		out = string(sign) + out
	}
	return out
}

// isAnagramByReverse is 5 times slower then isAnagram
func isAnagramByReverse(s, s1 string) bool {
	n := len(s)

	if n != len(s1) {
		return false
	}
	return string(revers([]rune(s1))) == s
}

func revers(slice []rune) []rune {

	n := len(slice)
	reversed := make([]rune, n)
	for i := n - 1; i >= 0; i-- {
		reversed[n-1-i] = slice[i]
	}
	return reversed

	// with interface{} as param and return
	// switch t := slice.(type) {
	// case []rune:
	// 	n := len(t)
	// 	reversed := make([]rune, n)
	// 	for i := n - 1; i >= 0; i-- {
	// 		reversed[n-1-i] = t[i]
	// 	}
	// 	return reversed

	// }
	// return slice
}

func isAnagram(s, s1 string) bool {
	n := len(s)

	if n != len(s1) {
		return false
	}

	r := []rune(s)
	rn := len(r)
	var j int
	for _, v := range s1 {
		if r[rn-1-j] != v {
			return false
		}
		j++
	}

	return true
}
