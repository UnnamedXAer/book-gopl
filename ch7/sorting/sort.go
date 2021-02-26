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
	order string
	// @thought: its not the correct way of doing for http use because it shares the state
	// @thought: across all clients therefore all of the clients will get the same order
	// @thought: and override others order if decide to reorder on their page,
	// @thought: but can be use to simple CLI program.
	fieldsOrder []string
	reverse     bool
}

// NewStatefulSort returns pointer to new statefulSort object
func NewStatefulSort(tracks []*Track, fields []string) *StatefulSort {
	x := &StatefulSort{
		t:     tracks,
		order: "asc",
		fieldsOrder: []string{
			"Title",
			"Artist",
			"Album",
			"Year",
			"Length",
			// will hold temporarily previous field when updating sort order
			"",
		},
		reverse: false,
	}
	for _, f := range fields {
		x.SetSortBy(f, "asc")
	}
	return x
}

// GetOrder returns sorting order and the fields order
func (x StatefulSort) GetOrder() (string, []string) {
	return x.order, x.fieldsOrder[:len(x.fieldsOrder)-1]
}

// SetSortBy sets the primary sort field of the entity in slice and the order
// of the primary column sort
func (x *StatefulSort) SetSortBy(field string, order string) error {
	if field != "Title" &&
		field != "Artist" &&
		field != "Album" &&
		field != "Year" &&
		field != "Length" {
		return fmt.Errorf("invalid field name %q", field)
	}
	if order != "" && order != "asc" && order != "desc" {
		return fmt.Errorf(`invalid order value %q, acceptable "asc", "desc" or ""`, order)
	}

	x.order = order

	for i := len(x.fieldsOrder) - 1; i >= 1; i-- {
		x.fieldsOrder[i] = x.fieldsOrder[i-1]
	}
	x.fieldsOrder[0] = field

	var found bool
	for i := 1; i < len(x.fieldsOrder)-1; i++ {
		if found == false && x.fieldsOrder[i] == field {
			found = true
		}
		if found {
			x.fieldsOrder[i] = x.fieldsOrder[i+1]
		}
	}
	return nil
}

func (x StatefulSort) Len() int {
	fmt.Print(x.reverse, x.fieldsOrder, "\n")
	return len(x.t)
}
func (x StatefulSort) Less(i, j int) bool {

	if x.reverse {
		return less(x.t[j], x.t[i], x.fieldsOrder, x.order)
	}
	return less(x.t[i], x.t[j], x.fieldsOrder, x.order)
}
func (x StatefulSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }

func less(x, y *Track, fieldsOrder []string, order string) bool {
	var i int
	if order == "desc" {
		l, ok := lessField(y, x, fieldsOrder[i])
		if ok {
			return l
		}
		i++
	}
	for ; i < len(fieldsOrder)-1; i++ {
		l, ok := lessField(x, y, fieldsOrder[i])
		if ok {
			return l
		}
	}
	return false
}

func lessField(x, y *Track, field string) (less bool, ok bool) {
	switch field {
	case "Title":
		if x.Title != y.Title {
			return x.Title < y.Title, true
		}
	case "Artist":
		if x.Artist != y.Artist {
			return x.Artist < y.Artist, true
		}
	case "Album":
		if x.Album != y.Album {
			return x.Album < y.Album, true
		}
	case "Year":
		if x.Year != y.Year {
			return x.Year < y.Year, true
		}
	case "Length":
		if x.Length != y.Length {
			return x.Length < y.Length, true
		}
	default:
		return false, false
	}
	return false, false
}
