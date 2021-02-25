package main

import "fmt"

type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

type StatefulSort struct {
	t     []*Track
	order []string
}

// NewStatefulSort returns pointer to new statefulSort object
func NewStatefulSort(tracks []*Track, fields []string) *StatefulSort {
	x := &StatefulSort{
		t: tracks,
		order: []string{
			"Title",
			"Artist",
			"Album",
			"Year",
			"Length",
			// will hold temporarily previous field when updating sort order
			"",
		},
	}
	for _, f := range fields {
		x.SetSortBy(f)
	}
	return x
}

// SetSortBy sets the primary sort field of the entity in slice
func (x StatefulSort) SetSortBy(field string) {
	if field != "Title" &&
		field != "Artist" &&
		field != "Album" &&
		field != "Year" &&
		field != "Length" {
		panic(fmt.Errorf("invalid field name %q", field))
	}
	// x.order = append(x.order, field)
	for i := len(x.order) - 1; i >= 1; i-- {
		x.order[i] = x.order[i-1]
	}
	x.order[0] = field

	var found bool
	for i := 1; i < len(x.order)-1; i++ {
		if found == false && x.order[i] == field {
			found = true
		}
		if found {
			x.order[i] = x.order[i+1]
		}
	}
	return
}

func (x StatefulSort) Len() int {
	fmt.Print(x.order, "\n")
	return len(x.t)
}
func (x StatefulSort) Less(i, j int) bool {

	return less(x.t[i], x.t[j], x.order)
}
func (x StatefulSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }

func less(x, y *Track, order []string) bool {
	for i := 0; i < len(order)-1; i++ {
		switch order[i] {
		case "Title":
			if x.Title != y.Title {
				return x.Title < y.Title
			}
		case "Artist":
			if x.Artist != y.Artist {
				return x.Artist < y.Artist
			}
		case "Album":
			if x.Album != y.Album {
				return x.Album < y.Album
			}
		case "Year":
			if x.Year != y.Year {
				return x.Year < y.Year
			}
		case "Length":
			if x.Length != y.Length {
				return x.Length < y.Length
			}
		default:
			return false
		}
	}
	return false
}
