// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
//
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/net/html"
)

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		// fmt.Println("return cancelled true")
		return true
	default:
		return false
	}
}

func main() {
	doWork(os.Args[1:])
}

func doWork(cmdArgs []string) {

	// debugging
	// const successPanicMsg = "success panic"
	// defer func() {
	// 	result := recover()
	// 	if result != nil {
	// 		x, ok := result.(string)
	// 		if ok && x == successPanicMsg {
	// 			fmt.Println("WE DID WELL!")
	// 			return
	// 		}
	// 		panic(result)
	// 	}
	// }()

	fmt.Println()
	// lists of URLs, may have duplicates
	worklist := make(chan []string)
	unseenLinks := make(chan string) // de-duplicated URLs

	// Start with the command-line arguments.
	go func() {
		if cancelled() {
			return
		}
		worklist <- cmdArgs
	}()

	go func() {
		c := 0
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if c > 0 {
				// debbuggin on block to see running goroutines
				// second enter will trigger the panic
				panic("keyboard triggered panic")

			}
			fmt.Println("`````````````````````````Closed chan 'done'")
			close(done)
			c++
		}
		// os.Stdin.Read(make([]byte, 1))
	}()

	wgRoot := sync.WaitGroup{}
	wgCrawl := sync.WaitGroup{}
	n := 200
	wgRoot.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("@@@@@1 recovered here from ", r)
					// panic(r)
				}
			}()
			for link := range unseenLinks {
				wgCrawl.Add(1)
				foundLinks := crawl(link)
				go func(foundLinks []string) {
					select {
					case worklist <- foundLinks:
					case <-done:
					}
					wgCrawl.Done()
				}(foundLinks)
			}
			wgRoot.Done()
		}()
	}

	go readWorklist(worklist, unseenLinks)

	<-done
	// clean up
	fmt.Println("CLEAN UP------------ START")
	// fmt.Printf("waiting for extracts %d to finish, worklist #%d\n", extracts, len(worklist))
	wgCrawl.Wait()

	// fmt.Println("about to close unseenlinks")
	close(unseenLinks)
	// fmt.Println("about to drain unseenlinks", len(unseenLinks))
	for range unseenLinks {
	}

	wgRoot.Wait()
	// fmt.Println(">>>about to close worklist", len(worklist))
	close(worklist)
	// fmt.Println(">>>about to drain worklist", len(worklist))
	for range worklist {

	}
	// fmt.Println(">>>worklist drained", len(worklist))
	// fmt.Println("CLEAN UP------------ END")

	fmt.Println("remaining extracts: ", extracts, "-", extractsDone, " = ", extracts-extractsDone)
	// fmt.Println("all work done")

	// debugging
	// panic(successPanicMsg)
}

func readWorklist(worklist <-chan []string, unseenLinks chan<- string) {
	// defer func() { fmt.Println("leaving the readWorklist()") }()

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)

	for {
		select {
		case list, ok := <-worklist:
			if !ok {
				log.Println("breaking the worklist 'loop' by  !OK")
				return
			}
			fmt.Println("read from worklist", len(list))
			func(list []string) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("@@@@@2 recovered here from ", r)
						// panic(r)
					}
				}()
				for _, link := range list {
					if !seen[link] {
						seen[link] = true
						unseenLinks <- link
					}
				}
			}(list)
		case <-done:
			// log.Println("breaking the worklist 'loop'")
			return
		}
	}
}

//
var cnt int64

func crawl(url string) []string {
	cnt++
	fmt.Println(cnt, "", url)
	list, err := extract(url)
	if err != nil {
		log.Println(err, list)
		// log.Println("crawl err")
		return nil
	}
	return list
}

// these are not necessary, they are here for debugging purposes
// while learning to cancel goroutines
var extracts int = 0
var extractsDone int = 0
var mu = sync.Mutex{}
var mu2 = sync.Mutex{}

func extract(url string) ([]string, error) {
	mu.Lock()
	extracts++
	mu.Unlock()
	token := make(chan struct{})
	defer func() {
		token <- struct{}{}
	}()
	ctx, reqCancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-token:
		case <-done:
			reqCancel()
			<-token //cancell req and release the token
		}
		mu2.Lock()
		extractsDone++
		// fmt.Println("extract done #", extractsDone)
		mu2.Unlock()
	}()

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
