package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/unnamedxaer/book-gopl/ch5/myhtml"
	"golang.org/x/net/html"
)

func main() {
	err := title("https://jsonplaceholder.typicode.com/")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Great!!!")
}

func title(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	t := res.Header.Get("Content-Type")
	if strings.HasPrefix(t, "text/html") == false {
		return fmt.Errorf("content of %q is %q, wanted %q", url, t, "text/html**")
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("parsing %q as HTML: %v", url, err)
	}

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				fmt.Println(n.FirstChild.Data)
			}
		}

	}
	myhtml.ForEachNode(doc, visitNode, nil)

	return nil
}
