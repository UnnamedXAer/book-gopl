package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	memo "github.com/unnamedxaer/book-gopl/ch9/ex9_3_memo"
)

func httpGetBody(url string, done memo.Done) (interface{}, error) {
	token := make(memo.Done)
	ctx, cancelReq := context.WithCancel(context.Background())
	defer func() {
		token <- struct{}{}
	}()

	go func() {
		select {
		case <-done:
			cancelReq()
			<-token
		case <-token:
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request fail: %s", res.Status)
	}

	return ioutil.ReadAll(res.Body)
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

func cancelFunc(done *memo.Done, key string) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		close(*done)
	}
}

func main() {
	var m = memo.New(HTTPGetBody)

	var done memo.Done = make(memo.Done)
	output := []string{}
	go cancelFunc(&done, "")

	for url := range incomingURLs() {
		done = make(memo.Done)

		start := time.Now()
		value, err := m.Get(url, done)
		if err != nil {
			log.Println(err)
		}
		var n int
		if value != nil {
			n = len(value.([]byte))
		}
		output = append(output, fmt.Sprintf("% 40s, % 15s, % 8d, bytes\n",
			url, time.Since(start), n))

		select {
		// I'm not sure if this closing is necessary.
		// Can anyone tell?
		case <-done:
		default:
			close(done)
		}
	}

	fmt.Println()
	fmt.Println(output)
}
