package views

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/unnamedxaer/book-gopl/ch4/ghweb/viewutil"
)

var layoutsDir = "web/layouts"
var flashRotator int = 0

type View struct {
	Template *template.Template
	Layout   string
}

type ViewData struct {
	Flashes    map[string]string
	RenderTime viewutil.ViewTime
	Data       interface{}
}

// NewView parses and returns new view with given layout
func NewView(layout string, funcs template.FuncMap, files ...string) *View {
	files = append(layoutFiles(), files...)
	t, err := template.New("").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Render sends view data to user. Render is deprecated, use Execute instead
func (v *View) Render(w http.ResponseWriter, data interface{}) error {

	vd := ViewData{
		Flashes: flashes(),
		Data:    data,
	}

	return v.Template.ExecuteTemplate(w, v.Layout, vd)
}

// Execute executes template and write it's results to w
//
//  b := &bytes.Buffer{}
//  err = issue.Execute(b, v)
//  if err != nil {
// 		return newAppError(
// 			err,
// 			"Sorry, because of some errors we couldn't print your page.",
// 			http.StatusInternalServerError)
// 	}
// 	fmt.Fprint(w, b)
func (v *View) Execute(w io.Writer, data interface{}) error {

	vd := ViewData{
		Flashes: flashes(),
		Data:    data,
	}
	return v.Template.ExecuteTemplate(w, v.Layout, vd)
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

func flashes() map[string]string {
	// return map[string]string{
	// 	"warning": "You are about to exceed your plan limits!",
	// }

	flashRotator = flashRotator + 1
	if flashRotator%3 == 0 {
		return map[string]string{
			"warning": "You are about to exceed your plan limits!",
		}
	}

	return map[string]string{}
}
