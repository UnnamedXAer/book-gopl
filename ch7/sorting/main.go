package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("33m8s")},
	{"Go 2", "Moby 2", "Moby 2", 1992, length("3m37s")},
	{"Go 2", "Moby 2", "Moby", 1991, length("3m37s")},
	{"Go 2", "Moby", "Moby", 1992, length("3m35s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go", "Moby", "Moby", 1990, length("3m37s")},
	{"Go", "Moby", "Moby", 1990, length("3m38s")},
	{"Go", "Moby", "Moby", 1990, length("3m33s")},
	{"Go", "Moby", "Moby", 1990, length("3m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
}

var sortableTracks = NewStatefulSort(tracks, nil)

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func printTracks(tracks []*Track, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)
	// tw := os.Stdout
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "-----", "-----", "-----", "-----")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}

	tw.Flush() // calculate column widths and print table
}

func printTracksHTML(w http.ResponseWriter, tracks []*Track, fieldsOrder []string, order string) {
	headers := []string{"Title", "Artist", "Album", "Year", "Length"}

	if order == "" || order == "asc" {
		order = "desc"
	} else {
		order = "asc"
	}

	data := map[string]interface{}{
		"Headers":      headers,
		"Tracks":       tracks,
		"OrderedField": fieldsOrder[0],
		"Order":        order,
	}

	b := new(bytes.Buffer)
	err := htmlTpl.Execute(b, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("\n,%s / %v\n\n", order, fieldsOrder)))
	w.Write(b.Bytes())
}

func main() {
	http.HandleFunc("/", tracksHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.SetOutput(os.Stdout)
	log.Println("Sever up and runing. http://localhost:3030")
	fmt.Fprintln(os.Stderr, http.ListenAndServe(":3030", nil))
}

func tracksHandler(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	order := r.URL.Query().Get("order")

	if field != "" {
		sortableTracks.SetSortBy(field, order)
		sort.Sort(sortableTracks)
	}
	_, fieldsOrder := sortableTracks.GetOrder()
	printTracksHTML(w, tracks, fieldsOrder, order)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	return
}

// func main1() {
// 	printTracks(tracks, os.Stdout)
// 	fmt.Println(sort.IsSorted(byArtist(tracks)))
// 	sort.Sort(byArtist(tracks))
// 	fmt.Println(sort.IsSorted(byArtist(tracks)))
// 	sort.Sort(customSort{tracks, func(x, y *Track) bool {
// 		if x.Title != y.Title {
// 			return x.Title < y.Title
// 		}
// 		if x.Year != y.Year {
// 			return x.Year < y.Year
// 		}
// 		if x.Length != y.Length {
// 			return x.Length < y.Length
// 		}
// 		return false
// 	}})
// 	printTracks(tracks, os.Stdout)
// 	ss := NewStatefulSort(tracks, nil)
// 	fmt.Println()
// 	sort.Sort(ss)
// 	printTracks(tracks, os.Stdout)

// 	for _, field := range os.Args[1:] {
// 		fmt.Println()
// 		ss.SetSortBy(field, "")
// 		sort.Sort(ss)
// 		printTracks(tracks, os.Stdout)
// 	}
// }
