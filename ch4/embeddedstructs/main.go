package main

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	X, Y int
}

type Circle struct {
	Point
	Radius int
}

type Wheel struct {
	Circle
	Spokes int
}

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies []Movie = []Movie{
	{Title: "Too fast", Year: 2000, Color: true,
		Actors: []string{"Fast guy 1", "Faster guy 2"},
	},
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Cool Hand Luke", "Humphrey Bogart", "Ingrid Bergman"},
	},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steven McQueen", "Jacquelie Bisset"},
	},
}

func main() {

	b, err := json.MarshalIndent(movies, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	var m []Movie

	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
}

func theWeel() {
	var w Wheel
	w.Y = 123
	w.Circle.Point.X = 12
	w.Radius = 10
	w.Spokes = 30

	fmt.Printf("%#v", w)
}
