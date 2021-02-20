package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/unnamedxaer/book-gopl/ch4/ghweb/viewutil"
)

type GithubUserResults struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []*GitHubUser
	FetchedAt         viewutil.ViewTime
}

type GitHubUser struct {
	Login   string
	HTMLURL string `json:"html_url"`
	Type    string
	NodeID  string `json:"node_id"`
}

func userSearchHandler(w http.ResponseWriter, r *http.Request) *appError {
	errs := map[string]string{}

	n, found := getNameParam(r.URL.Query())
	qerr := r.URL.Query().Get("err")

	if qerr != "" {
		errs["requestErr"] = qerr
	}

	if found && n == "" {
		errs["name"] = "Name is required"
	}

	v := Data{
		PageTitle: "Contact",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData: map[string]map[string]string{
			"errors": errs,
		},
	}
	b := &bytes.Buffer{}
	err := searchUserView.Execute(b, &v)
	if err != nil {
		return newAppError(
			err,
			"Sorry, because of some errors we couldn't print your page.",
			http.StatusInternalServerError)
	}

	fmt.Fprint(w, b)
	return nil
}

func userSearchResultsHandler(w http.ResponseWriter, r *http.Request) *appError {
	n, found := getNameParam(r.URL.Query())

	if found && n == "" {
		// mb http.Redirect to /search-user?n=
		return userSearchHandler(w, r)
	}

	ur, err := fetchUsersByName(n)
	if err != nil {
		r.URL.Query().Add("err", err.Error())
		// return userSearchHandler(w, r)
		// http.Redirect(w, r, "/search-user?err="+url.QueryEscape(err.Error()), http.StatusInternalServerError)
		return newAppError(err, err.Error(), http.StatusInternalServerError)
	}

	v := Data{
		PageTitle: "Contact",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  ur,
	}
	ur.FetchedAt = viewutil.NewTime()
	//render
	b := &bytes.Buffer{}

	err = usersView.Execute(b, v)
	if err != nil {
		return newAppError(
			err,
			"Sorry, because of some errors we couldn't print your page.",
			http.StatusInternalServerError)
	}
	fmt.Fprint(w, b)
	return nil
}

func fetchUsersByName(n string) (*GithubUserResults, error) {
	q := url.QueryEscape(fmt.Sprintf("%s in:login %[1]s in:name", n))
	url := "https://api.github.com/search/users?q=" + q
	l.Println("fetch user with url:", url)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	ur := &GithubUserResults{}

	err = json.NewDecoder(res.Body).Decode(ur)
	if err != nil {
		return nil, err
	}

	l.Printf("%v", ur)

	return ur, nil
}

func getNameParam(q url.Values) (n string, found bool) {
	for k, v := range q {
		if k == "n" {
			n = v[0]
			found = true
			return
		}
	}
	return
}
