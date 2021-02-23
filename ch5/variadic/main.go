package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func main() {
	// fmt.Println(MaxUint, bits.UintSize, MaxInt, MinInt)
	// println(max())
	// println(max(5, 6, 3, 6, 3, 7, 4, 2, 1, 0))
	// println(min())
	// println(min(5, 6, 3, 6, 3, 7, 4))

	// println()

	// println(max1(8))
	// println(max1(5, 6, 3, 6, 3, 7, 4, 2, 1, 0))
	// println()
	// println()

	// println(join("X"))
	// println(join("1", "2", "x", "__", "--", "||"))
	// strings.Join("1")

	res, _ := http.Get("https://golang.org/pkg/net/http/")
	defer res.Body.Close()
	doc, _ := html.Parse(res.Body)

	images := ElementsByTagName(doc, "img")
	headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")

	fmt.Println(len(images))
	fmt.Println(len(headings))
}

func max1(arg int, args ...int) int {
	count := len(args)
	if count == 0 {
		return arg
	}
	tmp := arg
	for i := 0; i < count; i++ {
		if tmp < args[i] {
			tmp = args[i]
		}
	}

	return tmp
}

func max(args ...int) int {
	count := len(args)
	if count == 0 {
		return MaxInt
	}
	tmp := args[0]
	for i := 1; i < count; i++ {
		if tmp < args[i] {
			tmp = args[i]
		}
	}

	return tmp
}

func min(args ...int) int {
	count := len(args)
	if count == 0 {
		return MinInt
	}
	tmp := args[0]
	for i := 1; i < count; i++ {
		if tmp > args[i] {
			tmp = args[i]
		}
	}

	return tmp
}

func join(sep string, s ...string) (j string) {
	count := len(s)
	if count == 0 {
		return
	}

	for i := 0; i < count-1; i++ {
		j += s[i] + sep
	}

	j += s[count-1]

	return
}

func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {

	// if contains(name, doc.Data) {
	// 	l = append(l, doc)
	// }

	return checkNode(doc, name...)
}

func checkNode(n *html.Node, name ...string) []*html.Node {
	var l []*html.Node
	if n.Type == html.ElementNode && contains(name, n.Data) {
		l = append(l, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		l = append(l, checkNode(c, name...)...)
	}

	return l
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if str == v {
			return true
		}
	}

	return false
}
