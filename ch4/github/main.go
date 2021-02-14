package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResults struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResults, error) {
	q := url.QueryEscape(strings.Join(terms, " "))

	res, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", res.Status)
	}

	var results IssuesSearchResults
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return &results, nil
}

func main() {
	issueses, err := SearchIssues([]string{"is:open", "json", "decoder"})
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(issueses, "", "  ")

	fmt.Printf("%s", b)
}
