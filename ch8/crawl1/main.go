package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unnamedxaer/book-gopl/ch5/links"
)

func main() {
	worklist := make(chan []string)

	// Start with the command-line arguments.
	go func() {
		worklist <- os.Args[1:]
	}()

	seen := make(map[string]bool)
	for list := range worklist {
		for _, url := range list {
			if !seen[url] {
				seen[url] = true
				go func(url string) {
					worklist <- crawl(url)
				}(url)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}
