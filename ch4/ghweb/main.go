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
	issue      *views.View
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
	issue = views.NewView("bootstrap", "web/views/issues.html")

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", staticHandler)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/issues", issuesHandler)

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
		responseOn500Error(w, err)
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
		responseOn500Error(w, err)
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
		responseOn500Error(w, err)
	}
}

func issuesHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, r.URL.Path)
	// err := r.ParseForm()
	// if err != nil {
	// 	responseOn500Error(w, err)
	// }

	un := r.FormValue("username")
	rn := r.FormValue("reponame")
	var ir interface{}
	var err error
	if un != "" && rn != "" {
		ir, err = fetchIssuesByUserRepo(un, rn)
	} else {
		k := r.Form["keywords"]
		ir, err = fetchIssuesByKeywords(k)
	}
	if err != nil {
		responseOn500Error(w, err)
		return
	}

	v := Data{
		PageTitle: "Contact",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  ir,
	}

	err = issue.Render(w, v)
	if err != nil {
		responseOn500Error(w, err)
		return
	}

}
