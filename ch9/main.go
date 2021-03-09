package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	memo "github.com/unnamedxaer/book-gopl/ch9/memo1"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func incomingURLs() []string {
	input := bufio.NewScanner(os.Stdin)
	s := []string{}
	for input.Scan() {
		s = append(s, input.Text())
	}

	return s
}

func main() {

	m := memo.New(httpGetBody)
	for _, url := range incomingURLs() {
		start := time.Now()
		value, err := m.Get("http://google.com")
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s, %s, %d, bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}
