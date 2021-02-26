package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var s = sync.Mutex{}

type database map[string]dollar
type dollar float32

func (d dollar) String() string {
	return fmt.Sprintf("$%.2f", d)
}

func (db database) list(w http.ResponseWriter, r *http.Request) {

	tpl := template.New("table")
	tpl.Parse(`<html><body><Table><tr><th>Name</th>Price<th></th></tr>{{range $name, $price := .}}
	<tr><td>{{$name}}</td><td>{{$price}}</td></tr>
	{{end}}</Table></body></html>`)

	// data := make(map[string]interface{}, len(db))

	err := tpl.Execute(w, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func (db database) add(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")

	s.Lock()
	if _, ok := (db)[item]; ok == true {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "item %q already exists\n", item)
		return
	}

	p64, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not a valid price %q\n", price)
		return
	}
	(db)[item] = dollar(p64)
	log.Println(db)
	fmt.Fprintf(w, "%s", (db))
	s.Unlock()
}

func (db database) update(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")

	s.Lock()
	if _, ok := (db)[item]; ok == false {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item %q\n", item)
		return
	}

	p64, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not a valid price %q\n", price)
		return
	}
	(db)[item] = dollar(p64)
	log.Println(db)
	fmt.Fprintf(w, "%s", (db)[item])
	s.Unlock()
}

func (db database) delete(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	s.Lock()
	delete(db, item)
	s.Unlock()

	http.NoBody.WriteTo(w)
}

func main() {
	db := database{
		"shoes":      123.2,
		"knee wraps": 199.99,
	}

	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/add", db.add)
	mux.HandleFunc("/delete", db.delete)

	go func() {
		http.HandleFunc("/list", db.list)
		http.HandleFunc("/price", db.price)
		http.HandleFunc("/update", db.update)
		http.HandleFunc("/add", db.add)
		http.HandleFunc("/delete", db.delete)
		log.Println("Server up na runing on 'localhost:3031'")
		log.Fatalln(http.ListenAndServe("localhost:3031", nil))
	}()

	log.Println("Server up na runing on 'localhost:3030'")
	log.Fatalln(http.ListenAndServe("localhost:3030", mux))
}
