package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {

	f, _ := os.Open("../outline3/d.html")

	doc, _ := html.Parse(f)

	s, err := soleTitle(doc)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s)
}

// soleTitle returns the text of the first non-empty title element
// in doc, and an error if there was not exactly one.
func soleTitle(doc *html.Node) (title string, err error) {

	type bailout struct{}

	defer func() {
		switch p := recover(); p {
		case nil:
			// no panic
		case bailout{}:
			// "expected" panic
			err = fmt.Errorf("multiple title elements")
		default:
			// unexpected panic; carry on panicking
			panic(p)
		}
	}()

	// Bail out recurion if we find more then one non-empty title.
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			if title != "" {
				panic(bailout{})
			}
			title = n.FirstChild.Data
		}
	}, nil)

	if title == "" {
		return "", fmt.Errorf("no title element 😔")
	}
	return title, nil
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
