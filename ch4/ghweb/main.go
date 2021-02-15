package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":3030", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	x, err := R
	fmt.Fprint(w, "hello home")
}
