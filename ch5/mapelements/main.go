package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mapelements: %/n", err)
		os.Exit(1)
	}

	// for _, link := range visit(nil, doc) {
	// 	fmt.Println(link)
	// }
	// m := map[string]int{}
	// m = visit(m, doc)
	// fmt.Printf("\ndifferent elements used: %v\n", len(m))

	// for n, c := range m {
	// 	fmt.Println(n, ":", c)
	// }
	m := []string{}
	m = visitText(m, doc)

	for n, c := range m {
		fmt.Println(n, ":", c)
	}
}

func visit(m map[string]int, n *html.Node) map[string]int {
	if n.Type == html.ElementNode {
		m[n.Data]++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		m = visit(m, c)
	}

	return m
}

func visitText(m []string, n *html.Node) []string {
	if n.Data == "script" || n.Data == "style" || n.Data == "noscript" {
		return m
	}
	if n.Type == html.TextNode {
		m = append(m, n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		m = visitText(m, c)
	}

	return m
}
