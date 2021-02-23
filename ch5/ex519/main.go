package main

import "fmt"

func main() {
	x := r()
	fmt.Println(x)
}

func r() (w string) {
	defer func() {
		p := recover()
		w = p.(string)
	}()
	panic("w")
}
