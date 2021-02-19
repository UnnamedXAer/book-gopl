package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/unnamedxaer/book-gopl/ch4/ghweb/web/views"
)

var (
	index      *views.View
	contact    *views.View
	issues     *views.View
	issue      *views.View
	errView    *views.View
	layoutsDir = "web/layouts"
	l          *log.Logger
)

var cnt int

// Data represents data passed to template
type Data struct {
	PageTitle string
	Author    string
	UserName  string
	AppName   string
	ViewData  interface{}
}

func main() {
	l = log.New(os.Stdout, "> ", log.LstdFlags)

	// http.HandleFunc("/favicon.ico", faviconHandler)
	index = views.NewView("bootstrap", "web/views/index.html")
	contact = views.NewView("bootstrap", "web/views/contact.html")
	issues = views.NewView("bootstrap", "web/views/issues.html")
	issue = views.NewView("bootstrap", "web/views/issue.html")
	errView = views.NewView("bootstrap", "web/views/error.html")

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", staticHandler)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/issues", issuesHandler)
	http.HandleFunc("/issue", issueHandler)
	http.HandleFunc("/error", errorHandler)

	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
		l.Println("/favicon.ico", "Referer:", r.Header["Referer"])
		rw.Header().Set("Content-Type", "image/x-icon")
		rw.Write([]byte{})
	})

	l.Println("Server available on http://localhost:3030")
	err := http.ListenAndServe(":3030", nil)
	checkErr(err)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./web/assets/favicon.ico")
	if err != nil {
		responseOn500Error(w, r, err)
		return
	}

	fmt.Fprint(w, b)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, "/", r.URL.Path)

	v := Data{
		PageTitle: "Home",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
	}

	err := index.Render(w, v)
	if err != nil {
		responseOn500Error(w, r, err)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, r.URL.Path)

	v := Data{
		PageTitle: "Contact",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
	}

	err := contact.Render(w, v)
	if err != nil {
		responseOn500Error(w, r, err)
	}
}

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

func issueHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, r.URL.Path)

	ids, ok := r.URL.Query()["id"]

	if ok == false || len(ids) == 0 {
		http.Error(w, fmt.Sprint("missing the issue node_id ('id=<string>' - query param)"), http.StatusBadRequest)
		return
	}
	nodeID := ids[0]

	var err error
	if nodeID == "" {
		http.Error(w, fmt.Sprint("missing the issue node_id value ('id=<string>' - query param)"), http.StatusBadRequest)
		return
	}
	iss, err := getIssue(nodeID)

	if err != nil {
		responseOn500Error(w, r, err)
		return
	}

	v := Data{
		PageTitle: "Issue",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  iss,
	}

	err = issue.Render(w, v)
	if err != nil {
		responseOn500Error(w, r, err)
		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, r.URL.Path)

	errText := r.URL.Query().Get("t")

	if errText != "" {
		errText += "\n"
	}
	errText = "Please try again later."

	v := Data{
		PageTitle: "Issue",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  errText,
	}

	err := issue.Render(w, v)
	if err != nil {
		responseOn500Error(w, r, err)
		return
	}
}
