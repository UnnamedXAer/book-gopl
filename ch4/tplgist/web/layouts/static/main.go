package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	layoutsDir = "web/layouts"
	viewsDir   = "web/views"
	contactTpl *template.Template
	searchTpl  *template.Template
)

func main() {
	funcs := template.FuncMap{
		"addOne": func(x int) int { return x + 1 },
	}

	contactTpl = makeTemplate(funcs, viewsDir+"/contact.html")
	searchTpl = makeTemplate(nil, viewsDir+"/search.html")

	http.HandleFunc("/contact", contactHandler)

	log.Println("Server available on http://localhost:" + "3030")
	err := http.ListenAndServe(":"+"3030", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	b := bytes.Buffer{}

	// I do not call tpl.Execute directly to http.ResponseWriter because
	// in case of error the error message we want to show would
	// be appended to rendered template (the part before error occurs)
	err := contactTpl.Execute(&b, map[string]map[string]string{
		"errors": {
			"mail":  "email is required",
			"title": "title max length is 50 characters",
		},
	})
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, &b)
}

func makeTemplate(funcs template.FuncMap, files ...string) *template.Template {
	files = append(layoutFiles(), files...)
	t, err := template.New("").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return t
}

func layoutFiles() []string {
	var files []string
	err := filepath.Walk(layoutsDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return err
	})

	if err != nil {
		panic(err)
	}
	return files
}
