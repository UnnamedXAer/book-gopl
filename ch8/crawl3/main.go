package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unnamedxaer/book-gopl/ch5/links"
)

func main() {
	fmt.Println()
	// lists of URLs, may have duplicates
	worklist := make(chan []string)
	unseenLinks := make(chan string) // de-duplicated URLs

	// Start with the command-line arguments.
	go func() {
		worklist <- os.Args[1:]
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}
	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	n := 1
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				n++
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
	fmt.Println("done main")
}

var cnt int64

func crawl(url string) []string {
	cnt++
	fmt.Println(cnt, "", url)
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}
