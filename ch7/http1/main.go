package main

import (
	"fmt"
	"log"
	"net/http"
)

type database map[string]dollar
type dollar float32

func (d dollar) String() string {
	return fmt.Sprintf("$%.2f", d)
}

func (db database) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// b, err := json.Marshal(&db)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	// w.Write(append(b, '\n', '\n'))

	switch r.URL.Path {
	case "/list":
		for item, price := range db {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	case "/price":
		item := r.URL.Query().Get("item")
		price, ok := db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no such item %q\n", item)
			return
		}
		fmt.Fprintf(w, "%s", price)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such page: %s\n", r.URL)
	}
	return
}

func main() {
	db := database{
		"shoes":      123.2,
		"knee wraps": 199.99,
	}
	log.Fatalln(http.ListenAndServe("localhost:3030", db))
}
