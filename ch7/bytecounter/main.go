package bytecounter

import (
	"bufio"
	"bytes"
	"io"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	r := bytes.NewReader(p)

	scanner := bufio.NewScanner(r)
	var n int
	for scanner.Scan() {
		n++
		// fmt.Println(n, scanner.Text())
	}
	*c += LineCounter(n)
	return n, nil
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	r := bytes.NewReader(p)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var n int
	for scanner.Scan() {
		n++
		// fmt.Println(n, scanner.Text())
	}
	*c += WordCounter(n)
	return n, nil
}

type XWriter struct {
	n int64
	w io.Writer
}

func (c *XWriter) Write(p []byte) (int, error) {
	c.n += int64(len(p))
	return len(p), nil
}
func (c *XWriter) N() int64 {
	return c.n
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	xw := XWriter{
		w: w,
	}
	return &xw, &xw.n
}

// func main() {

// 	f, err := os.Open("main.go")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer f.Close()
// 	b, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	s := string(b)

// 	// var c ByteCounter

// 	// n, err := fmt.Fprintln(&c, s)
// 	// if err != nil {
// 	// 	fmt.Fprintf(os.Stderr, "fail while writing, written bytes %d, error: %v", n, err)
// 	// }
// 	// fmt.Println("n", n)
// 	// fmt.Println("c", c)
// 	// fmt.Println()

// 	// var lc LineCounter

// 	// n, err = fmt.Fprintln(&lc, s)
// 	// if err != nil {
// 	// 	fmt.Fprintf(os.Stderr, "fail while writing, written lines %d, error: %v", n, err)
// 	// }
// 	// fmt.Println("n", n)
// 	// fmt.Println("lc", lc)

// 	// var wc WordCounter

// 	// n, err = fmt.Fprintln(&wc, s)
// 	// if err != nil {
// 	// 	fmt.Fprintf(os.Stderr, "fail while writing, written words %d, error: %v", n, err)
// 	// }
// 	// fmt.Println("n", n)
// 	// fmt.Println("wc", wc)
// 	var c ByteCounter
// 	var tmpn int
// 	xw, n := CountingWriter(&c)

// 	fmt.Println()
// 	fmt.Println("*n", *n)
// 	tmpn, err = xw.Write([]byte(s))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Println()
// 	fmt.Println("*n", *n)
// 	fmt.Println("tmpn", tmpn)
// 	tmpn, err = xw.Write([]byte(s))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Println()
// 	fmt.Println("*n", *n)
// 	fmt.Println("tmpn", tmpn)
// 	fmt.Println("xw.n", xw.(*XWriter).N())
// }
