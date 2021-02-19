package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	htmlTemplate "html/template"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/unnamedxaer/book-gopl/ch4/ghweb/viewutil"
)

const IssuesByKeywordsURL = "https://api.github.com/search/issues"
const IssuesByUserRepoURLTplText = "https://api.github.com/repos/{{.Username}}/{{.Reponame}}/issues"

// const IssueByNum = "https://api.github.com/repos/{{.Username}}/{{.Reponame}}/issues{{if .Num}}/{{.Num}}{{end}}"

type IssuesSearchResults struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
	FetchedAt  viewutil.ViewTime
}

type IssuesSearchResultsUserRepo struct {
	TotalCount int
	RepoName   string
	UserName   string
	Items      []*Issue
	FetchedAt  viewutil.ViewTime
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	URL       string
	Title     string
	State     string
	User      *User
	CreatedAt viewutil.ViewTime `json:"created_at"`
	Body      htmlTemplate.HTML // in Markdown format
	NodeID    string            `json:"node_id"`
}

var (
	issuesByUserRepoURLTpl *template.Template
	issueCacheWithKeywords map[string]*IssuesSearchResults         = map[string]*IssuesSearchResults{}
	issueCacheWithUserRepo map[string]*IssuesSearchResultsUserRepo = map[string]*IssuesSearchResultsUserRepo{}
)

func init() {
	issuesByUserRepoURLTpl = template.Must(template.New("url").Parse(IssuesByUserRepoURLTplText))
}

func fetchIssuesByKeywords(params []string) (*IssuesSearchResults, error) {
	fmt.Println(issueCacheWithKeywords)
	q := url.QueryEscape(strings.Join(params, " "))

	cached := issueCacheWithKeywords[q]
	if cached != nil {
		return cached, nil
	}
	res, err := http.Get(IssuesByKeywordsURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	ir := IssuesSearchResults{}

	err = json.NewDecoder(res.Body).Decode(&ir)
	if err != nil {
		return nil, err
	}
	ir.FetchedAt = viewutil.NewIssueTime()
	issueCacheWithKeywords[q] = &ir

	for _, v := range ir.Items {
		b := markdown.ToHTML([]byte(v.Body), nil, nil)
		v.Body = htmlTemplate.HTML(b)
	}

	return &ir, nil
}

func fetchIssuesByUserRepo(un, rn string) (*IssuesSearchResultsUserRepo, error) {
	fmt.Println(issueCacheWithUserRepo)
	qUn := url.QueryEscape(un)
	qRn := url.QueryEscape(rn)
	b := bytes.Buffer{}
	err := issuesByUserRepoURLTpl.Execute(&b, struct {
		Username string
		Reponame string
	}{
		qUn,
		qRn,
	})
	if err != nil {
		return nil, err
	}

	url := b.String()
	l.Println(url)

	cachedWithUserRepo := issueCacheWithUserRepo[url]
	if cachedWithUserRepo != nil {
		return cachedWithUserRepo, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ir := IssuesSearchResultsUserRepo{}
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("Not Exists")
	}

	err = json.NewDecoder(res.Body).Decode(&ir.Items)
	if err != nil {
		return nil, err
	}
	ir.FetchedAt = viewutil.NewIssueTime()
	ir.TotalCount = len(ir.Items)
	ir.RepoName = rn
	ir.UserName = un

	for _, v := range ir.Items {
		b := markdown.ToHTML([]byte(v.Body), nil, nil)
		v.Body = htmlTemplate.HTML(b)
	}

	issueCacheWithUserRepo[url] = &ir
	return &ir, nil
}

func findIssue(list []*Issue, nodeID string) *Issue {
	for _, item := range list {
		if item.NodeID == nodeID {
			return item
		}
	}
	return nil
}

func getIssue(nodeID string) (*Issue, error) {
	var cashedIssue interface{}

	for _, v := range issueCacheWithUserRepo {
		cashedIssue = findIssue(v.Items, nodeID)
		if cashedIssue != nil {
			break
		}
	}
	if cashedIssue == nil {
		for _, v := range issueCacheWithKeywords {
			cashedIssue = findIssue(v.Items, nodeID)
			if cashedIssue != nil {
				break
			}
		}
	}

	if cashedIssue == nil {
		return nil, fmt.Errorf("not_found")
	}

	issue, ok := cashedIssue.(*Issue)
	if ok {
		t := issue.CreatedAt.Time()
		if t.After(time.Now().Add(-24 * time.Hour)) {
			return issue, nil
		}
	}

	var err error
	issue, err = fetchIssue(issue.URL)
	if err != nil {
		// we will return the cached issue
	}
	return issue, nil
}

func fetchIssue(url string) (*Issue, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	iss := Issue{}
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("not_found")
	}

	err = json.NewDecoder(res.Body).Decode(&iss)
	if err != nil {
		return nil, err
	}

	b := markdown.ToHTML([]byte(iss.Body), nil, nil)
	iss.Body = htmlTemplate.HTML(b)

	return &iss, nil
}
