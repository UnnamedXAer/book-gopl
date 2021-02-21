// Package links provides a link-extraction function.
package links

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// func main() {
// 	l, err := Extract(os.Args[1])
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	fmt.Printf(strings.Join(l, "\n"))
// }

//Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and resturns the links in the HTML document.
func Extract(url string) ([]string, error) {
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
