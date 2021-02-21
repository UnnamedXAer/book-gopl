package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	// fmt.Println(os.Args)
	// f, _ := os.Open("doc.html")
	// fmt.Println(string(b))
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v", err)
	}

	outline(nil, doc)
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data)
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// stack = append(stack, c.Data)

		outline(stack, c)
	}
}
