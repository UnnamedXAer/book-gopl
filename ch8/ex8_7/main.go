package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/unnamedxaer/book-gopl/ch5/links"
)

var seen = map[string]bool{}
var seenLock = sync.Mutex{}

func main() {
	baseLink := flag.String("link", "", "base page link")
	flag.Parse()
	if *baseLink == "" {
		log.Fatalln("the 'link' flag is required")
	}

	fetchAndSave(*baseLink)
	fmt.Println("done main")
}

func fetchAndSave(baseLink string) {
	seen = map[string]bool{}
	fmt.Println("\nSTART: ", baseLink)

	url, err := url.Parse(baseLink)
	if err != nil {
		log.Fatalln(err)
	}
	baseDomain := url.Host

	dir, err := prepareDir(baseDomain)
	if err != nil {
		log.Fatalln(err)
	}

	docs := make(chan [2]string, 100)
	worklist := make(chan []string, 100)

	worklist <- []string{baseLink}

	go func() {
		// go processWorklist(worklist, docs, baseDomain)
		processWorklist(worklist, docs, baseDomain)
		close(docs)
	}()

	n := 100
	wgSaveDocs := sync.WaitGroup{}
	wgSaveDocs.Add(n)
	for i := 0; i < n; i++ {
		go saveDocs(&wgSaveDocs, docs, baseDomain, dir)
	}
	wgSaveDocs.Wait()
}

func processWorklist(worklist chan []string, docs chan<- [2]string, baseDomain string) {
	n := len(worklist)

	for ; n > 0; n-- {
		fmt.Println("n", n)
		items := <-worklist
		// for items := range worklist {
		for _, link := range items {
			seenLock.Lock()
			if seen[link] {
				continue
			}

			seen[link] = true
			seenLock.Unlock()

			u, err := url.Parse(link)
			if err != nil {
				log.Printf("url parse: %v\n", err)
				continue
			}
			if u.Host != baseDomain {
				log.Printf("different hosts, skipped: %q\n", link)
				continue
			}
			n++
			// yes it will result in double get request
			// we should extract links from fetch results.
			// and do not call links.Extract
			go getPageLinks(worklist, link)
			content, err := fetchPage(link)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
				continue
			}

			docs <- [2]string{link, content}
		}
	}
	// close(docs)
}

func getPageLinks(wl chan<- []string, link string) {
	log.Printf("crawling: %s", link)
	links, err := crawl(link)
	if err != nil {
		log.Println(err)
	}
	log.Println()

	wl <- links
}

func saveDocs(wg *sync.WaitGroup, docs <-chan [2]string, baseDomain string, dir string) {
	for doc := range docs {
		path := filepath.Join(dir, pathifyURL(doc[0])+".go.html")
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()

		log.Printf("saving: %s\n", path)
		if _, err = f.WriteString(doc[1]); err != nil {
			log.Println(err)
		}
	}

	wg.Done()
}

func fetchPage(link string) (content string, err error) {
	log.Printf("feetching: %s\n", link)
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("fetching %q  resulted in %q", link, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func pathifyURL(u string) string {
	return strings.ReplaceAll(url.PathEscape(u), ":", "")
}

func crawl(url string) ([]string, error) {
	pageLinks, err := links.Extract(url)
	if err != nil {
		return nil, err
	}
	return pageLinks, nil
}

func prepareDir(baseDomain string) (string, error) {
	dir := filepath.Join("html", pathifyURL(baseDomain))

	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return "", fmt.Errorf("mkdir: %v", err)
	}
	return dir, nil
}
