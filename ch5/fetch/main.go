package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	fn, n, err := Fetch("https://stackoverflow.com/questions/10485743/contains-method-for-a-slice")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%d bytes wrote to file %q", n, fn)
}

// Fetch downloads the URL and returns the
// name and length of the local file.
func Fetch(url string) (filename string, nbytes int64, err error) {
	timeout := time.Second * 2
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	filename = path.Base(res.Request.URL.Path)

	if filename == "/" {
		filename = "index.html"
	}

	f, err := os.Create(filename)
	if err != nil {
		return
	}

	defer func() {
		closeErr := f.Close()
		if err == nil {
			fmt.Println("file closed, close err", closeErr)
			err = closeErr
		}
	}()

	// b, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return
	// }
	// n, err := f.Write(b)
	// nbytes = int64(n)
	nbytes, err = io.Copy(f, res.Body)
	// if closeErr := f.Close(); err == nil {
	// 	err = closeErr
	// }
	return
}
