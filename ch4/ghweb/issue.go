package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

const IssuesByKeywordsURL = "https://api.github.com/search/issues"
const IssuesByUserRepoURLTplText = "https://api.github.com/repos/{{.Username}}/{{.Reponame}}/issues"

type issueTime time.Time

func NewIssueTime() issueTime {
	return issueTime(time.Now())
}

func (it issueTime) String() string {
	return time.Time(it).Format(time.RFC1123)
}

func (it *issueTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02T15:04:05Z", s)
	*it = issueTime(t)
	return
}

func (it *issueTime) MarshalJSON() ([]byte, error) {
	return []byte(it.String()), nil
}

type IssuesSearchResults struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
	FetchedAt  issueTime
}

type IssuesSearchResultsUserRepo struct {
	TotalCount int
	RepoName   string
	UserName   string
	Items      []*Issue
	FetchedAt  issueTime
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt issueTime `json:"created_at"`
	Body      string    // in Markdown format
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
	ir.FetchedAt = NewIssueTime()
	issueCacheWithKeywords[q] = &ir
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
	ir.FetchedAt = NewIssueTime()
	ir.TotalCount = len(ir.Items)
	ir.RepoName = rn
	ir.UserName = un

	issueCacheWithUserRepo[url] = &ir
	return &ir, nil
}
