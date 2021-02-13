package main

import "fmt"

func main() {
	// x := make([]int, 3, 5)
	// i := 0

	// x[0] = 2
	// x[1] = 3
	// x[2] = 4

	// y := x[:4]

	// fmt.Printf("% 2d, len=% 3d, cap=% 3d, %v\n", i, len(x), cap(x), x)
	// fmt.Printf("% 2d, len=% 3d, cap=% 3d, %v\n", i, len(y), cap(y), y)

	x := []int{}
	for i := 0; i < 10; i++ {
		x = appendInt(x, i)
		fmt.Printf("% 2d, len=% 3d, cap=% 3d, %v\n", i, len(x), cap(x), x)
	}
	x = []int{}
	fmt.Println()
	fmt.Println()
	for i := 0; i < 60; i += 3 {
		x = appendInt(x, i)
		fmt.Printf("% 2d, len=% 3d, cap=% 3d, %v\n", i, len(x), cap(x), x)
	}
}

func appendInt(x []int, y int) []int {
	// var z []int

	zlen := len(x) + 1

	if zlen <= cap(x) {
		z := x[:zlen]
		z[zlen-1] = y
		return z
	}

	c := 2 * cap(x)
	if zlen > c {
		c = zlen
	}

	z := make([]int, zlen, c)
	copy(z, x)
	z[zlen-1] = y
	return z
}

func appendIntBookVer(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		z = x[:zlen]
	} else {
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	z[len(x)] = y
	return z
}
