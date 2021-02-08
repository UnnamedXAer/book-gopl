package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		url = fixURL(url)
		fetchAndPrint(url)
	}
}

func fixURL(url string) string {
	if strings.HasPrefix(url, "http://") == false && strings.HasPrefix(url, "https://") == false {
		url = "http://" + url
	}
	return url
}

func fetchAndPrint(url string) {
	data, code, err := fetch2(url)
	fmt.Println()
	fmt.Println("url: " + url)
	fmt.Println("status code: " + strconv.Itoa(code))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't fetch from %q,\nstatus code: %d, error: %q", url, code, err)
		return
	}
	if data != "" {
		fmt.Println("data: " + data)
	}
}

func fetch(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", res.StatusCode, err
	}
	defer res.Body.Close()
	return string(b), res.StatusCode, nil
}

func fetch2(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		return "", res.StatusCode, err
	}
	defer res.Body.Close()
	return "", res.StatusCode, nil
}
