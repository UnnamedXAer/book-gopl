package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	// err = os.Remove(outputFileName)
	// if err != nil {
	// 	log.Printf("couldn't remove file: %q, err: %v", outputFileName, err)
	// }

	of, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer of.Close()
	// err = of.Truncate(0)
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = of.Seek(0, 0)
	// if err != nil {
	// 	panic(err)
	// }

	s := forEachNode(doc, startElement, endElement)

	n, err := of.WriteString(s)
	if err != nil {
		log.Println("write error: ", err)
	}
	log.Println(n, "byte were written")
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) string) string {
	var s string
	if pre != nil {
		s = pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s += forEachNode(c, pre, post)
	}

	if post != nil {
		s += post(n)
	}
	return s
}

func startElement(n *html.Node) (s string) {
	if n.Type == html.ElementNode {
		s = fmt.Sprintf("%*s<%s", depth*2, "", n.Data)
		for _, v := range n.Attr {
			s += fmt.Sprintf(` %s="%s"`, v.Key, v.Val)
		}

		if n.FirstChild == nil && (n.Data == "img") {
			s += " /"
		}
		s += ">\n"
		depth++
	} else if n.Type == html.TextNode {
		t := strings.ReplaceAll(strings.ReplaceAll(n.Data, "\t", ""), "\n", " ")
		if t != "" {
			s += fmt.Sprintf("%*s%s\n", depth*2, "", t)
		}
	}
	return
}

func endElement(n *html.Node) string {
	if n.Type == html.ElementNode {
		depth--
		if n.Data != "img" {
			return fmt.Sprintf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
	return ""
}
