package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/unnamedxaer/book-gopl/ch5/links"
)

var copyExists = map[string]string{}

func main() {
	baseWorklist := []string{
		"https://www.petefreitag.com/cheatsheets/ascii-codes/",
	}
	seen := map[string]bool{}

	for _, item := range baseWorklist {
		fmt.Printf("-% *s %s\n", 5, "", item)
		breadthFirst(seen, crawlWithCopy, item)
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item
func breadthFirst(seen map[string]bool, f func(baseHosts, pageURL string) []string, baseItem string) {
	worklist := []string{baseItem}
	u, err := url.Parse(baseItem)
	if err != nil {
		log.Printf("incorrect url skipped (%s), error: %v", baseItem, err)
		return
	}
	if u.String() != baseItem {
		fmt.Printf("the link is different then its parsed version\nlink:   %q\nparsed: %q\n", baseItem, u.String())
	}
	baseHost := u.Host

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if seen[item] == false {
				// host :=
				seen[item] = true
				worklist = append(worklist, f(baseHost, item)...)
			}
		}
	}
}

func crawlWithCopy(baseHosts, pageURL string) []string {
	terminateOnLimitExceeded(10)
	fmt.Println(pageURL)
	list, err := links.Extract(pageURL)
	if err != nil {
		log.Print(err)
	} else {
		for _, v := range list {
			err := makeCopy(baseHosts, v)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return list
}

func makeCopy(baseHost, urlToCopy string) error {
	utocopy, err := url.Parse(urlToCopy)

	if baseHost != utocopy.Host {
		return fmt.Errorf("hosts mismatch page host: %q, link: %q", baseHost, urlToCopy)
	}

	if copyExists[urlToCopy] != "" {
		return fmt.Errorf("copy of %q already exists", urlToCopy)
	}

	dir := "./web/" + pathifyURL(baseHost)

	err = os.MkdirAll(dir, 0644)
	if err != nil {
		return fmt.Errorf("unable to create directory, error: %v", err)
	}

	fn := pathifyURL(urlToCopy) + ".go.html"
	path := dir + "/" + fn

	// @todo: should create a file and if the file already exists set "copyExists" and return
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {

		if os.IsExist(err) {
			copyExists[urlToCopy] = path
			return fmt.Errorf("copy of %q already exists", urlToCopy)
		}
		return fmt.Errorf("unable to determine if copy of %q already exists, error: %v", urlToCopy, err)
	}

	return readAndSavePage(f, urlToCopy, path)
}

func pathifyURL(u string) string {
	return strings.ReplaceAll(url.PathEscape(u), ":", "")
}

func readAndSavePage(f *os.File, u, p string) error {
	res, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("unable to make copy of %q, error: %v", u, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("unable to make copy of %q, reason: %q", u, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to make copy of %q, cannot read body, error: %v", u, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("unable to write the copy of %q to a file, error: %v", u, err)
	}
	copyExists[u] = strconv.Itoa(len(copyExists))
	return nil
}

var cnt int

func terminateOnLimitExceeded(n int) {
	cnt++
	if n > 100 {
		n = 100
	}
	if cnt > n {
		fmt.Printf("limit (%d) exceeded\n", n)
		os.Exit(0)
	}
}
