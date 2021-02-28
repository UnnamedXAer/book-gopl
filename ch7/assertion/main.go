package main

import (
	"fmt"
	"io"
	"os"

	"github.com/unnamedxaer/book-gopl/ch7/bytecounter"
)

func main() {
	var w io.Writer
	w = os.Stderr
	f, ok := w.(*os.File)
	if ok == false {
		fmt.Printf("f false")
	}
	f.Write([]byte("emejzing"))

	w = new(bytecounter.ByteCounter)
	// b, ok := w.(io.ReadWriter)
	// if ok == false {
	// 	fmt.Printf("b false")
	// }
	w.Write([]byte{'A'})

	w.Write([]byte("Fupa"))
	rw := w.(io.ReadWriter)
	rw.Write([]byte("Pafu"))
	w = new(bytecounter.ByteCounter)
	w = rw
	w = rw.(io.Writer)

	fmt.Fprintf(rw, "\n\n\nXXX")
}
