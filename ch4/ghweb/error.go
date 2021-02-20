package main

import (
	"log"
	"net/http"

	"github.com/unnamedxaer/book-gopl/ch4/ghweb/config"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

func newAppError(err error, msg string, code int) *appError {
	return &appError{
		err,
		msg,
		code,
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func responseOn500Error(w http.ResponseWriter, r *http.Request, err error) {
	l.Println(err)
	// if f, ok := w.(http.Flusher); ok {
	// 	f.Flush()
	// 	l.Println("response flushed")
	// }
	// w.WriteHeader(http.StatusInternalServerError)
	// fmt.Fprintf(w, err.Error())
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// http.Redirect(w, r, "/error?t="+err.Error(), http.StatusSeeOther)
	if config.C.Debbug {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
