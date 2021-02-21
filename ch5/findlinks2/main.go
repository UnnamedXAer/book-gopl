package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	for _, arg := range os.Args[1:] {
		// links, err := fetchLinks(arg)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "findlinks2 %v\n", err)
		// }
		// for _, link := range links {
		// 	fmt.Println(link)
		// }

		w, i, err := CountWordsAndImages(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "findlinks2 - count words err: %v", err)
			continue
		}
		fmt.Printf("document from %q have %d words and %d images\n", arg, w, i)
	}
}

func fetchLinks(url string) ([]string, error) {
	res, err := http.Get(url)
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
	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		if n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		} else if n.Data == "img" || n.Data == "style" || n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					links = append(links, a.Val)
				} else if a.Key == "srcset" {
					links = append(links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}

// CountWordsAndImages makes HTTP GET request to given url for the HTML
// then counts words and images in the response
func CountWordsAndImages(url string) (words, images int, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("getting %q result in: %q", url, res.Status)
		return
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		err = fmt.Errorf("parsing %q as HTML: %v", url, err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n.Type == html.ElementNode {
		if n.Data == "img" {
			images++
		}
	} else if n.Type == html.TextNode {
		r := strings.NewReader(n.Data)
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			words++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w, i := countWordsAndImages(c)
		words = words + w
		images = images + i
	}

	return
}
