package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func issuesHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, r.URL.Path)

	un := r.FormValue("username")
	rn := r.FormValue("reponame")
	var ir interface{}
	var err error
	if un != "" && rn != "" {
		ir, err = fetchIssuesByUserRepo(un, rn)
	} else if len(r.Form["keywords"]) > 0 {
		k := r.Form["keywords"]
		ir, err = fetchIssuesByKeywords(k)
	} else {
		http.Error(w, fmt.Sprint("Missing query params"), http.StatusBadRequest)
		return
	}

	if err != nil {
		responseOn500Error(w, r, err)
		return
	}

	v := Data{
		PageTitle: "Contact",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  ir,
	}

	err = issues.Render(w, v)
	if err != nil {
		responseOn500Error(w, r, err)
		return
	}

}

func issueHandler(w http.ResponseWriter, r *http.Request) *appError {
	cnt++
	fmt.Println(cnt, r.URL.Path)

	ids, ok := r.URL.Query()["id"]

	if ok == false || len(ids) == 0 {
		// http.Error(w, fmt.Sprint("missing the issue node_id ('id=<string>' - query param)"), http.StatusBadRequest)
		return newAppError(
			fmt.Errorf("missing the issue node_id ('&id=<string>' - query param)"),
			"Your url is lack of issue ID.",
			http.StatusBadRequest)
	}
	nodeID := ids[0]

	var err error
	if nodeID == "" {
		// http.Error(w, fmt.Sprint("missing the issue node_id value ('id=<string>' - query param)"), http.StatusBadRequest)
		return newAppError(
			fmt.Errorf("missing value of the issue node_id value ('id=<string>' - query param)"),
			("Your url is lack of the issue ID value."),
			http.StatusBadRequest)
	}
	iss, err := getIssue(nodeID)

	if err != nil {
		// responseOn500Error(w, r, err)
		return newAppError(
			err,
			"Sorry, we were not able to get information about the issue.",
			http.StatusInternalServerError)
	}

	v := Data{
		PageTitle: "Issue",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  iss,
	}

	b := &bytes.Buffer{}
	err = issue.Execute(b, v)
	if err != nil {
		return newAppError(
			err,
			"Sorry, because of some errors we couldn't print your page.",
			http.StatusInternalServerError)
	}
	fmt.Fprint(w, b)
	return nil
}
