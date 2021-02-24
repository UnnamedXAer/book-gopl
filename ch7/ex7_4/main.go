package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
	s   string
	pos int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.pos >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.pos:])

	r.pos += int64(n)

	return n, nil
}

func NewReader(s string) *Reader {
	return &Reader{
		s,
		0,
	}
}

func main() {
	// fReads()
}

func fReads() {
	s := "Ann has a cat"
	// var x io.Reader
	b := make([]byte, 12, 30)
	rr := strings.NewReader(s)
	n, err := rr.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "11read %d bytes, error: %v\n", n, err)
	}
	fmt.Printf("11bytes read %d, text: %s\n", n, b)
	n, err = rr.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "12read %d bytes, error: %v\n", n, err)
	}
	fmt.Printf("12bytes read %d, text: %s\n", n, b)
	n, err = rr.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "13read %d bytes, error: %v\n", n, err)
	}
	fmt.Printf("13bytes read %d, text: %s\n", n, b)
	fmt.Println()
	b = make([]byte, 12, 30)
	rr2 := NewReader(s)
	///////////////////////////////////////////////////////////////
	n, err = rr2.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "1read %d bytes, error: %v\n", n, err)
	}
	fmt.Printf("21bytes read %d, text: %s\n", n, b)
	n, err = rr2.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "2read %d bytes, error: %v\n", n, err)
	}
	fmt.Printf("22bytes read %d, text: %s\n", n, b)
	n, err = rr2.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "3read %d bytes, error: %v\n", n, err)
	}
	fmt.Printf("23bytes read %d, text: %s\n", n, b)
}
