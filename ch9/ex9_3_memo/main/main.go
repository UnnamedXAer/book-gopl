package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	memo "github.com/unnamedxaer/book-gopl/ch9/ex9_3_memo"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var HTTPGetBody = httpGetBody

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {

		for _, url := range []string{
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
		} {
			ch <- url
		}
		close(ch)
	}()

	return ch
}

func cancel(done memo.Done, key string) {
	os.Stdin.Read(make([]byte, 1))
	fmt.Printf(">> closing done for the key: %s\n", key)
	close(done)
}

func main() {
	var m = memo.New(HTTPGetBody)

	for url := range incomingURLs() {
		fmt.Println()
		var done memo.Done = make(memo.Done)
		key := url
		go cancel(done, key)
		start := time.Now()
		value, err := m.Get(url, done)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("% 40s, % 15s, % 8d, bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}
