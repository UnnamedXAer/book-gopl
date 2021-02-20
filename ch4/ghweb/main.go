package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/unnamedxaer/book-gopl/ch4/ghweb/config"
	"github.com/unnamedxaer/book-gopl/ch4/ghweb/web/views"
)

var (
	index          *views.View
	contact        *views.View
	issues         *views.View
	issue          *views.View
	searchUserView *views.View
	usersView      *views.View
	errView        *views.View
	layoutsDir     = "web/layouts"
	l              *log.Logger
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

func init() {

}

func main() {
	l = log.New(os.Stdout, "> ", log.LstdFlags)

	// http.HandleFunc("/favicon.ico", faviconHandler)
	index = views.NewView("bootstrap", nil, "web/views/index.html")
	contact = views.NewView("bootstrap", nil, "web/views/contact.html")
	issues = views.NewView("bootstrap", nil, "web/views/issues.html")
	issue = views.NewView("bootstrap", nil, "web/views/issue.html")
	searchUserView = views.NewView("bootstrap", nil, "web/views/users-search.html")
	funcs := template.FuncMap{
		"addOne": func(x int) int { return x + 1 },
	}
	usersView = views.NewView("bootstrap", funcs, "web/views/users.html")
	errView = views.NewView("bootstrap", nil, "web/views/error.html")

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", staticHandler)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/issues", issuesHandler)
	// http.HandleFunc("/issue", issueHandler)
	http.Handle("/issue", appHandler(issueHandler))
	http.Handle("/search-user", appHandler(userSearchHandler))
	http.Handle("/users", appHandler(userSearchResultsHandler))
	http.HandleFunc("/error", errorHandler)

	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
		l.Println("/favicon.ico", "Referer:", r.Header["Referer"])
		rw.Header().Set("Content-Type", "image/x-icon")
		rw.Write([]byte{})
	})

	l.Println("Server available on http://localhost:" + strconv.Itoa(int(config.C.PORT)))
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.C.PORT)), nil)
	checkErr(err)
}
