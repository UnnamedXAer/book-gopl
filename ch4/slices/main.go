package main

import (
	"fmt"
	"unicode"
)

func main() {

	// s := []int{1, 2, 3, 4, 5, 6, 7}
	// idx := 4
	// fmt.Printf("s[idx:]: %v\n", s[idx:])
	// fmt.Printf("s[idx+1:]: %v\n", s[idx+1:])

	// // s1 := make([]int, len(s), len(s))
	// s1 := make([]int, 1, 1)
	// s2 := make([]int, 10, 10)
	// fmt.Printf("before: %v\n", s)
	// copy(s1, s[2:5])
	// copy(s2, s[2:5])

	// // s1 := []int{9, 10, 11}
	// // copy(s, s1)

	// fmt.Printf("s:  %v\n", s)
	// fmt.Printf("s1:  %v\n", s1)
	// fmt.Printf("s2:  %v\n", s2)

	// fmt.Println()
	// s = []int{1, 2, 3, 4, 5, 6, 7}
	// fmt.Printf("before: %v\n", s)
	// s = remove(s, 4)
	// fmt.Printf("after:  %v\n", s)
	// fmt.Println()
	// s = []int{1, 2, 3, 4, 5, 6, 7}
	// fmt.Printf("2. before: %v\n", s)
	// s, _ = removeWithErr(s, 4)
	// fmt.Printf("2. after:  %v\n", s)

	// arr := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// fmt.Printf("1. before: %v\n", arr)
	// // reverse(&arr)
	// // reverses(arr[:])
	// rotateLeft(arr[:], 3)
	// fmt.Printf("1. after: %v\n", arr)
	// arr = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// strs := []string{
	// 	"OS",
	// 	"ENV",
	// 	"ENV",
	// 	"ENV",
	// 	"GO",
	// 	"JS",
	// 	"JS",
	// 	"TS",
	// 	"OS",
	// 	"TMP",
	// 	"TMP",
	// 	"ENV",
	// 	"YA",
	// }
	b := []byte{'a', 'b', ' ', ' ', ' ', 'a', 'b', '\t', '\n', '\v', '\f', '\r', ' ', '1', 0x09, 0x0C, '1', 'a', 'b', ' ', '\n', 'a', 'b'}
	fmt.Printf("2. before: %v\n", b)
	fmt.Printf("2. before: %s\n", string(b))
	// reverse(&arr)
	// reverses(arr[:])
	// strs = rmAdjacentDuplicates(strs)
	b = squashSpaces(b)
	fmt.Printf("2. after: %v\n", b)
	fmt.Printf("2. after: %s\n", string(b))
}

func squashSpaces(b []byte) []byte {
	s := string(b)
	r := []rune(s)
	count := len(r)

	inRow := false
	for i := 0; i < count; i++ {
		fmt.Println(r[i])
		if unicode.IsSpace(r[i]) && unicode.IsSpace(r[i-1]) {
			if inRow == false {
				inRow = true
				r[i] = ' '
			} else {
				copy(r[i-1:], r[i:])
			}
			count--
			i--
			continue
		}
		inRow = false
	}
	r = r[:count]
	s = string(r)
	b = []byte(s)
	return b
}

func rmAdjacentDuplicates(s []string) []string {
	count := len(s)
	for i := 1; i < count; i++ {
		if s[i-1] == s[i] {
			copy(s[i-1:], s[i:])
			count--
			i--
		}
	}
	s = s[:count]
	return s
}

func rotateLeft(s []int, n int) {
	reverses(s[:n])
	reverses(s[n:])
	reverses(s)
}

func reverse(a *[10]int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		// aj := a[j]
		// a[j] = a[i]
		// a[i] = aj
		a[i], a[j] = a[j], a[i]
	}
}
func reverses(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		sj := s[j]
		s[j] = s[i]
		s[i] = sj
	}
}

func remove(s []int, idx int) []int {
	copy(s[idx:], s[idx+1:])
	return s[:len(s)-1]
}

func remove2(s []int, idx int) []int {
	count := len(s)
	for i := idx + 1; i < count; i++ {
		s[i-1] = s[i]
	}
	return s[:count-1]
}

func removeWithErr(s []int, idx int) ([]int, error) {
	count := len(s)
	if idx >= count {
		return s, fmt.Errorf("cannot remove element of index %d, index out of range [%d]", idx, count-1)
	}
	if idx < 0 {
		return s, fmt.Errorf("cannot remove element at negative index %d", idx)
	}
	for i := idx + 1; i < count; i++ {
		s[i-1] = s[i]
	}

	return s[:count-1], nil
}
