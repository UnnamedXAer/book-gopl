package main

import "fmt"

func main() {
	in, out := pipline(5 * 1000000)
	in <- uint64(0)
	fmt.Println(<-out)
}

func pipline(n uint64) (in, out chan uint64) {
	out = make(chan uint64)
	first := out
	for i := uint64(0); i < n; i++ {
		in = out
		out = make(chan uint64)
		go func(in, out chan uint64) {
			x := <-in
			out <- x + 1
			close(out)
		}(in, out)
	}
	return first, out
}
