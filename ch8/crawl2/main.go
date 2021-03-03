package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unnamedxaer/book-gopl/ch5/links"
)

// tokens is a counting semaphore used to
// enforce a limit of 20 concurent requests
var tokens = make(chan struct{}, 20)

func main() {
	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	os.Exit(1)
	// }()

	worklist := make(chan []string)
	var n int // nu,ber of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() {
		worklist <- os.Args[1:]
	}()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		log.Println("n: ", n)
		list := <-worklist
		for _, url := range list {
			if !seen[url] {
				seen[url] = true
				n++
				go func(url string) {
					worklist <- crawlWithTokens(url)
				}(url)
			}
		}
	}
}

func crawlWithTokens(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire token
	list, err := links.Extract(url)
	<-tokens // release token
	if err != nil {
		log.Println(err)
	}
	return list
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}
