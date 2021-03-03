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

func main() {
	baseLink := flag.String("link", "", "base page link")
	flag.Parse()
	// fmt.Println("START: ", *baseLink)

	fetchAndSave(*baseLink)
}

func prepareDir(baseDomain string) (string, error) {
	dir := filepath.Join("html", pathifyURL(baseDomain))
	// fmt.Println(dir)

	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return "", fmt.Errorf("mkdir: %v", err)
	}
	return dir, nil
}

func fetchAndSave(baseLink string) {
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
	go func() {
		worklist <- []string{baseLink}
		close(worklist)
	}()

	go func() {
		go processWorklist(worklist, docs)
		processWorklist(worklist, docs)
		close(docs)

	}()

	n := 2
	wgSaveDocs := sync.WaitGroup{}
	wgSaveDocs.Add(n)
	for i := 0; i < n; i++ {
		go saveDocs(&wgSaveDocs, docs, baseDomain, dir)

	}
	wgSaveDocs.Wait()
}

func processWorklist(worklist <-chan []string, docs chan<- [2]string) {
	for items := range worklist {
		for _, link := range items {
			content, err := fetchPage(link)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch: %v", err)
				continue
			}

			docs <- [2]string{link, content}
		}
	}
	// close(docs)
}

func saveDocs(wg *sync.WaitGroup, docs chan [2]string, baseDomain string, dir string) {
	for doc := range docs {
		path := filepath.Join(dir, pathifyURL(doc[0])+".go.html")
		// log.Println(path)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		// fmt.Println()
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()

		if _, err = f.WriteString(doc[1]); err != nil {
			log.Println(err)
		}
	}

	wg.Done()
}

func fetchPage(link string) (content string, err error) {
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
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
