package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if appErr := fn(w, r); appErr != nil {
		l.Println(appErr.Error)
		w.WriteHeader(appErr.Code)
		v := Data{
			ViewData: appErr.Message,
		}
		if err := errView.Render(w, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Sorry, something broke here.")
		}
	}
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

	if r.URL.Path != "/" {
		fmt.Fprintln(w, "The aren't any resources for your URL.")
		return
	}

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

func errorHandler(w http.ResponseWriter, r *http.Request) {
	cnt++
	fmt.Println(cnt, r.URL.Path)

	errText := r.URL.Query().Get("t")

	if errText != "" {
		errText += "\n"
	}
	errText += "Please try again later."

	v := Data{
		PageTitle: "Issue",
		Author:    "Me",
		UserName:  "UnnamedXAer",
		AppName:   "Github Data",
		ViewData:  errText,
	}

	err := errView.Render(w, v)
	if err != nil {
		responseOn500Error(w, r, err)
		return
	}
}
