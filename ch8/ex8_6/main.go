package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/unnamedxaer/book-gopl/ch5/links"
)

type ldepth struct {
	link  string
	depth int
}

func main() {
	depth := flag.Int("depth", -1, "max crawl depth")
	flag.Parse()
	fmt.Println("Depth: ", *depth)
	// lists of URLs, may have duplicates
	worklist := make(chan []ldepth)
	// worklist := make(chan []string)
	unseenLinks := make(chan ldepth) // de-duplicated URLs
	// unseenLinks := make(chan string) // de-duplicated URLs

	// Start with the command-line arguments.
	go func(depth int) {
		// a better way would be to loop through the args and push urls
		i := 3
		if depth == -1 {
			// worklist <- os.Args[1:]
			i = 1
		}
		links := make([]ldepth, 0, len(os.Args[i:]))
		for _, l := range os.Args[i:] {
			if strings.HasPrefix(l, "http") == false {
				continue
			}
			links = append(links, ldepth{
				link:  l,
				depth: 1,
			})
		}
		worklist <- links
		// worklist <- os.Args[3:]
	}(*depth)

	if *depth == -1 {
		*depth = 1
	}

	linkInfoChan := make(chan string, 20)
	go func() {
		f, err := os.OpenFile("t1.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		for s := range linkInfoChan {
			if _, err = f.WriteString(s); err != nil {
				log.Println(err)
			}
		}
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for linkInfo := range unseenLinks {
				cnt++
				fmt.Printf("\rcnt: %d, depth: %d", cnt, linkInfo.depth)
				s := fmt.Sprintf("% 4d. depth: %d URL: %s\n", cnt, linkInfo.depth, linkInfo.link)
				linkInfoChan <- s
				foundLinks := crawl(linkInfo.link)
				// if linkInfo.depth == *depth {
				// 	continue
				// }
				go func(foundLinks []string, depth int) {
					links := make([]ldepth, len(foundLinks))
					for i, l := range foundLinks {
						links[i] = ldepth{
							link:  l,
							depth: depth + 1,
						}
					}
					worklist <- links
				}(foundLinks, linkInfo.depth)
			}

			close(linkInfoChan)
			log.Println("closed links chan")
		}()
	}
	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	n := 1
	for ; n > 0; n-- {
		list := <-worklist
		for _, linkInfo := range list {
			if linkInfo.depth > *depth {
				continue
			}

			if !seen[linkInfo.link] {
				n++
				seen[linkInfo.link] = true
				unseenLinks <- linkInfo
			}
		}
	}
	fmt.Println("done main")
}

var cnt int64

func crawl(url string) []string {
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}
