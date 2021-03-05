// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var cancel = make(chan struct{})

func cancelled() bool {
	select {
	case <-cancel:
		return true
	default:
		return false
	}
}

func main() {
	fmt.Println()
	// lists of URLs, may have duplicates
	worklist := make(chan []string)
	unseenLinks := make(chan string) // de-duplicated URLs

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(cancel)
	}()

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

	if cancelled() {
		return nil
	}

	cnt++
	fmt.Println(cnt, "", url)
	list, err := extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}

func extract(url string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	select {
	case <-cancel:
		<-req.Cancel
		fmt.Printf("cancelled: %s", req.URL.String())
	default:
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("getting %q result in: %q", url, res.Status)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				// parses URL as relative to the base URL
				link, err := res.Request.URL.Parse(a.Val)
				if err != nil {
					continue // skip bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}
