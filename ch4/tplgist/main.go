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
	issueTpl   *template.Template
)

// config/config.go
// var Dev = os.Getenv("ENV") == "development"
// mock
var config = struct {
	Dev bool
}{
	true || os.Getenv("ENV") == "development",
}

func main() {
	funcs := template.FuncMap{
		"addOne": func(x int) int { return x + 1 },
	}

	contactTpl = newTemplate(funcs, viewsDir+"/contact.html")
	issueTpl = newTemplate(nil, viewsDir+"/issue.html")

	http.HandleFunc("/contact", contactHandler)

	log.Println("Server available on http://localhost:" + "3030")
	err := http.ListenAndServe(":"+"3030", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func newTemplate(funcs template.FuncMap, files ...string) *template.Template {
	files = append(layoutFiles(), files...)
	tplName := "" // or filepath.Base(files[len(files)-1])
	t, err := template.New(tplName).Funcs(funcs).ParseFiles(files...)
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

func contactHandler(w http.ResponseWriter, r *http.Request) {
	// validate params
	// if (param == "" || ...) {...}
	errors := map[string]map[string]string{
		"errors": {
			"mail":    "email is required",
			"message": "message max length is 500 characters",
		},
	}
	// I do not call tpl.Execute directly to http.ResponseWriter because
	// in case of error the error message we want to show would
	// be appended to rendered template (the part before error occurs)
	b := bytes.Buffer{}
	err := contactTpl.ExecuteTemplate(&b, "bootstrap", errors)
	if err != nil {
		errText := http.StatusText(http.StatusInternalServerError)
		if config.Dev {
			errText = err.Error()
		}
		// or render error page
		http.Error(
			w,
			errText,
			http.StatusInternalServerError)
		return
	}

	//send mail ect. and render some page

	fmt.Fprint(w, &b)

}
