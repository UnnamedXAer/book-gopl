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
	defer log.Println("httpGetBody - end", done, url)
	token := make(memo.Done)
	ctx, cancelReq := context.WithCancel(context.Background())
	defer func() {
		token <- struct{}{}
	}()

	go func() {
		select {
		case <-done:
			log.Printf("cancelling request for %q", url)
			cancelReq()
			<-token
		case <-token:
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	// close(done)
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
		log.Printf(">> closing done, triggered by user %v, %s\n", *done, key)
		close(*done)
	}
	log.Println("@@@ leaving cancelFunc for", done, key)
}

func main() {
	var m = memo.New(HTTPGetBody)

	var done memo.Done = make(memo.Done)
	output := []string{}
	go cancelFunc(&done, "")

	for url := range incomingURLs() {
		done = make(memo.Done)
		fmt.Println()
		// log.Println("address of done:", done, url)

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
		case <-done:

		default:
			// log.Println("about to close done at the end of the loop iteration", done)
			close(done)
		}
	}

	fmt.Println()
	fmt.Println(output)
}
