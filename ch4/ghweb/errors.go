package main

import (
	"log"
	"net/http"
)

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
	http.Redirect(w, r, "/error?t="+err.Error(), http.StatusInternalServerError)
	// http.Error(w, err.Error(), http.StatusInternalServerError)
}
