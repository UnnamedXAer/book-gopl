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

func responseOn500Error(w http.ResponseWriter, err error) {
	l.Println(err)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
		l.Println("response flushed")
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
