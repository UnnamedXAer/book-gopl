package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

type pageInfo struct {
	idx int
	url,
	body string
}

func main() {
	fmt.Println()
	fmt.Println()
	wg := sync.WaitGroup{}
	pageInfo := mirroredQuery(&wg, os.Args[1:])

	lines := strings.Split(pageInfo.body, "\n")
	var fragment string
	for _, s := range lines {
		if s != "" && strings.ToLower(s) != "<!doctype html>" && len(lines) > 5 {
			fragment = s
			break
		}
	}

	if len(fragment) > 87 {
		fragment = fragment[:87] + "..."
	}
	fmt.Fprintf(os.Stdout, "\n  ````\n\nURL: %s\nidx: %d, num of bytes in body: %d\nText:\n%s\n\n \n  ````\n",
		pageInfo.url,
		pageInfo.idx,
		len(pageInfo.body),
		fragment)
	fmt.Println()
	wg.Wait()
}

func mirroredQuery(wg *sync.WaitGroup, links []string) *pageInfo {
	responses := make(chan *pageInfo)
	cancel := make(chan struct{})
	fmt.Println("Number od links:\n", len(links))
	wg.Add(len(links))
	for i := 0; i < len(links); i++ {
		go func(link string, i int) {
			defer func() {
				wg.Done()
			}()
			pageinfo, err := req(cancel, link)
			if err != nil {
				fmt.Fprintln(os.Stderr, "skipped, error: ", err)
				return
			}
			pageinfo.idx = i
			select {
			case responses <- pageinfo:
				close(cancel)

				fmt.Println("The first page: ", link)
			case <-cancel:
				fmt.Println("skipped, too late", link)
			}
		}(links[i], i)
	}
	return <-responses
}

func req(cancel chan struct{}, link string) (*pageInfo, error) {
	ctx, cancelReq := context.WithCancel(context.Background())
	token := make(chan struct{})
	defer func() {
		token <- struct{}{}
	}()
	go func() {
		select {
		case <-cancel:
			cancelReq()
			<-token
		case <-token:
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, fmt.Errorf("could create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request fail: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request fail: %s", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body of %q, %v", link, err)
	}

	return &pageInfo{
		0,
		link,
		string(b),
	}, nil
}
