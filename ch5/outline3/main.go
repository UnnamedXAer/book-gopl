package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

var depth int

const outputFileName = "r.html"

func main() {

	f, err := os.Open("d.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		log.Fatalln(err)
	}

	id := "footer-id"
	if len(os.Args) > 1 {
		id = os.Args[1]
	}

	n := ElementByID(doc, id)
	fmt.Printf("\nID: %s\nnode: %v", id, n)

}

func forEachNode(n *html.Node, id string, pre func(n *html.Node, id string) (found bool)) *html.Node {
	if pre != nil {
		if found := pre(n, id); found {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node := forEachNode(c, id, pre)
		if node != nil {
			return node
		}
	}

	return nil
}

func startElement(n *html.Node, id string) (found bool) {
	if n.Type == html.ElementNode {
		for _, v := range n.Attr {
			if v.Key == "id" && v.Val == id {
				return true
			}
		}
	}
	return false
}

// func endElement(n *html.Node) bool {
// 	if n.Type == html.ElementNode {
// 		depth--
// 		if n.Data != "img" {
// 			return fmt.Sprintf("%*s</%s>\n", depth*2, "", n.Data)
// 		}
// 	}
// 	return ""
// }

// ElementByID traversal the in given doc and returns first node with ID
func ElementByID(doc *html.Node, id string) *html.Node {
	return forEachNode(doc, id, startElement)
}
