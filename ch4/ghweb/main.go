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
	layoutsDir = "web/layouts"
	l          *log.Logger
)

var cnt int

// ViewData represents data passed to template
type ViewData struct {
	PageTitle    string
	Author       string
	UserName     string
	AppName      string
	ProjectTitle string
	Keywords     []string
}

func main() {
	l = log.New(os.Stdout, "> ", log.LstdFlags)

	// http.HandleFunc("/favicon.ico", faviconHandler)
	index = views.NewView("bootstrap", "web/views/index.html")
	contact = views.NewView("bootstrap", "web/views/contact.html")
	contact = views.NewView("bootstrap", "web/views/issues.html")

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", staticHandler)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/issues", issuesHandler)

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

	v := ViewData{
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

	v := ViewData{
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
	err := r.ParseForm()
	if err != nil {
		responseOn500Error(w, err)
	}
	k := r.Form["keywords"]

	v := ViewData{
		PageTitle:    "Contact",
		Author:       "Me",
		UserName:     "UnnamedXAer",
		AppName:      "Github Data",
		ProjectTitle: "FTS2020",
		Keywords:     k,
	}

	err = contact.Render(w, v)
	if err != nil {
		responseOn500Error(w, err)
	}
}
