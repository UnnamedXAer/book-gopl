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

var copyExists map[string]string

func main() {
	worklist := []string{
		"https://linuxize.com/post/what-does-chmod-777-mean/",
	}

	// breadthFirst(crawl, worklist)

	worklistWithHost := make(map[string]string, len(worklist))

	for _, v := range worklist {
		u, err := url.Parse(v)
		if err != nil {
			log.Printf("incorrect url skipped (%s), error: %v", v, err)
			continue
		}
		if u.String() != v {
			fmt.Printf("the link is different then it's parsed version\nlink:   %q\nparsed: %q\n", v, u.String())
		}
		worklistWithHost[v] = u.Host
	}
	fmt.Println(worklistWithHost)
	breadthFirst(worklistWithHost, crawlWithCopy, worklist)
	fmt.Println(strings.Join(worklist, "\n- "))
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item
func breadthFirst(worklistWithHost map[string]string, f func(worklistWithHost map[string]string, item string) []string, worklist []string) {
	seen := map[string]bool{}

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if seen[item] == false {
				// host :=
				seen[item] = true
				worklist = append(worklist, f(worklistWithHost, item)...)
			}
		}
	}
}

// func crawl(pageURL string) []string {
// 	terminateOnLimitExceeded(10)
// 	fmt.Println(pageURL)
// 	list, err := links.Extract(pageURL)
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	return list
// }

func crawlWithCopy(worklistWithHost map[string]string, pageURL string) []string {
	terminateOnLimitExceeded(10)
	fmt.Println(pageURL)
	list, err := links.Extract(pageURL)
	if err != nil {
		log.Print(err)
	} else {
		if host, ok := worklistWithHost[pageURL]; ok {
			for _, v := range list {
				makeCopy(host, v)
			}
		}
	}

	return list
}

func makeCopy(baseURLHost, urlToCopy string) {
	utocopy, err := url.Parse(urlToCopy)

	if baseURLHost != utocopy.Host {
		fmt.Printf("hosts mismatch page host: %q, link: %q\n", baseURLHost, urlToCopy)
		return
	}

	if copyExists[urlToCopy] != "" {
		fmt.Printf("copy of %q already exists\n", urlToCopy)
		return
	}

	dir := "./web/" + pathifyURL(baseURLHost)

	err = os.MkdirAll(dir, 0644)
	if err != nil {
		fmt.Printf("unable to create directory, error: %v\n", err)
		return
	}

	fn := pathifyURL(urlToCopy) + ".html"
	path := dir + "/" + fn

	// @todo: should create a file and if the file already exists set "copyExists" and return
	f, err := os.OpenFile(copyExists[urlToCopy], os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {

		if os.IsExist(err) {
			copyExists[urlToCopy] = path
			fmt.Printf("copy of %q already exists\n", urlToCopy)
			return
		}
		fmt.Println("-")
		// fmt.Printf("unable to determine if copy of %q already exists, error: %v\n", urlToCopy, err)
		return
	}

	readAndSavePage(f, urlToCopy, path)
}

func pathifyURL(u string) string {
	return strings.ReplaceAll(url.PathEscape(u), ":", "")
}

func readAndSavePage(f *os.File, u, p string) {
	res, err := http.Get(u)
	if err != nil {
		fmt.Printf("unable to make copy of %q, error: %v\n", u, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("unable to make copy of %q, reason: %q\n", u, res.Status)
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("unable to make copy of %q, cannot read body, error: %v\n", u, err)
		return
	}

	_, err = f.Write(b)
	if err != nil {
		fmt.Printf("unable to write the copy of %q to a file, error: %v\n", u, err)
		return
	}
	copyExists[u] = strconv.Itoa(len(copyExists))
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
