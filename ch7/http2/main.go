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

func (db database) list(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}
func (db database) price(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s", price)
}

func main() {
	db := database{
		"shoes":      123.2,
		"knee wraps": 199.99,
	}

	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.HandleFunc("/price", db.price)

	go func() {
		http.HandleFunc("/list", db.list)
		http.HandleFunc("/price", db.price)
		log.Println("Server up na runing on 'localhost:3031'")
		log.Fatalln(http.ListenAndServe("localhost:3031", nil))
	}()

	log.Println("Server up na runing on 'localhost:3030'")
	log.Fatalln(http.ListenAndServe("localhost:3030", mux))
}
