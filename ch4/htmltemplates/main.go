package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Item struct {
	Title   string
	DueDate time.Time
}

type ViewData struct {
	Name  string
	Items []Item
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

//https://www.alexedwards.net/blog/serving-static-sites-with-go
//https://gist.github.com/joyrexus/ff9be7a1c3769a84360f
func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveTemplate)

	log.Println("Server available on http://localhost:3030")
	err := http.ListenAndServe(":3030", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		log.Println("err:", err)
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmp, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println("err:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewData := ViewData{
		Name: "Dean",
		Items: []Item{
			{
				"clean room",
				time.Now().Add(time.Hour * 24),
			},
			{
				"wash car",
				time.Now().Add(time.Hour * 100),
			},
			{
				"pickup package",
				time.Now().Add(time.Hour * 48),
			},
		},
	}

	tmp.ExecuteTemplate(w, "layout", viewData)
	if err != nil {
		log.Println("err:", err)
		// http.StatusText(500)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
