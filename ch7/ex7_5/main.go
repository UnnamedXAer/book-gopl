package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
	reader io.Reader
	endn   int64
	idx    int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.idx >= r.endn {
		return 0, io.EOF
	}

	endIdx := len(p)
	if endIdx > int(r.endn) {
		endIdx = int(r.endn)
	}

	n, err = r.reader.Read(p[:endIdx])
	r.idx = int64(n)
	fmt.Println(r.endn, r.idx)
	return n, err
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &Reader{
		reader: r,
		endn:   n,
		idx:    0,
	}
}

var cnt int

func main() {
	s := "123456789_12" //3456789_123456789_123456789_123456789_"
	r := strings.NewReader(s)
	lr := LimitReader(r, 19)
	b := make([]byte, 3, 200)

	readAndPrint(lr, b)
	readAndPrint(lr, b)
	readAndPrint(lr, b)
	readAndPrint(lr, b)
	readAndPrint(lr, b)
	readAndPrint(lr, b)
}

func readAndPrint(lr io.Reader, b []byte) {
	cnt++
	b = make([]byte, len(b))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, cnt)
	n, err := lr.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on read, bytes read %d, error: %v\n", n, err)
	}
	fmt.Fprintf(os.Stdout, "bytes read %d, text: %s\n", n, b)
	// fmt.Fprintf(os.Stdout, "bytes %b\n", b)
}
