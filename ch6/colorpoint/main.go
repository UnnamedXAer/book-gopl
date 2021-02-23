package main

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/unnamedxaer/book-gopl/ch6/geometry"
)

type ColoredPoint struct {
	geometry.Point
	Color color.RGBA
}

func main() {

	var cp ColoredPoint
	cp.X = 100
	cp.Y = 200
	fmt.Println(cp)

	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	var p = ColoredPoint{geometry.Point{1, 1}, red}
	var q = ColoredPoint{geometry.Point{5, 4}, blue}

	p.ScaleBy(2)
	q.ScaleBy(8)
	fmt.Println(p, q)

	fmt.Println(Lookup("key"))
	fmt.Println(Lookup2("key"))

}

var (
	mu      sync.Mutex
	mapping = map[string]string{"key": "!234"}
)

func Lookup(key string) string {
	mu.Lock()
	v := mapping[key]
	mu.Unlock()
	return v
}

var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: map[string]string{"key": "!234$"},
}

func Lookup2(key string) string {
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}
