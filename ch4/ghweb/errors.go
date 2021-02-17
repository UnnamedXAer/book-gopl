package main

import (
	"fmt"
	"log"
	"net/http"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func responseOn500Error(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
