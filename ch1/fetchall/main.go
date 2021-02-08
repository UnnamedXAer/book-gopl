// fetch all url in paraller and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan string)

	for _, url := range os.Args[1:] {
		if strings.HasPrefix(url, "http") == false {
			fmt.Printf("wrong scheme, url %q\n", url)
			continue
		}
		fname := strings.Split(strings.Split(url, "//")[1], "/")[0]
		f, err := os.Create(fname + ".txt")
		if err != nil {
			fmt.Printf("create file: %q error: %v\n", fname, err)
			f = os.Stdout
		}
		defer f.Close()

		f2, err := os.Create(fname + "2" + ".txt")
		if err != nil {
			fmt.Printf("create file: %q error: %v\n", fname, err)
			f = os.Stdout
		}
		defer f2.Close()

		go fetch(url, f, 1, c)
		go fetch(url, f2, 2, c)
	}
	for range os.Args[1:] {
		println(<-c)
	}
	for range os.Args[1:] {
		println(<-c)
	}
	fmt.Printf("\n%.5fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, out io.Writer, n int, c chan<- string) {
	start := time.Now()

	res, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprintf("while reading %q, error: %v", url, err)
		return
	}

	nbytes, err := io.Copy(out, res.Body)
	defer res.Body.Close()
	if err != nil {
		c <- fmt.Sprintf("while reading %q, error: %v", url, err)
		return
	}

	durr := time.Since(start).Seconds()
	c <- fmt.Sprintf("%d) %.2fs elapsed, %8d bytes, %q", n, durr, nbytes, url)
}
