package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	g "github.com/unnamedxaer/book-gopl/ch6/geometry"
)

func main() {
	// f1()
	// f2()
	// f3()
	// time.Sleep(4 * time.Second)
	// fmt.Println("4s elapsed")
	// f4()
	f5()
}

func f5() {
	path := g.Path{{1, 2}, {4, 5}, {10, 20}, {-2, 1}}
	fmt.Println(path)
	path.TranslateBy(g.Point{5, 5}, true)
	fmt.Println(path)
	path.TranslateBy(g.Point{5, 0}, false)
	fmt.Println(path)

}

func f4() {
	p := g.Point{1, 2}
	q := g.Point{4, 6}

	distance := g.Point.Distance // method expression
	fmt.Println(distance(p, q))
	fmt.Printf("%T\n", distance)

	scale := (*g.Point).ScaleBy
	scale(&p, 2)
	fmt.Println(p)
	fmt.Printf("%T\n", scale)
}

type Rocket struct {
	Name string
}

func (r Rocket) Launch() {
	fmt.Printf("the rocket %v has been launched\n", r)
}

func f3() {
	r := Rocket{"perseverance"}
	timer := time.AfterFunc(1*time.Second, r.Launch)
	time.Sleep(1 * time.Second)

	timer.Reset(1 * time.Second)
	time.Sleep(1 * time.Second)
	timer.Reset(1 * time.Second)
	fmt.Println("timer reset", timer)
}

func f2() {
	p := g.Point{1, 2}
	q := g.Point{4, 6}
	distanceFomP := p.Distance // method value
	fmt.Println(distanceFomP(q), "==", p.Distance(q))
	// p.X = 100
	// fmt.Println(distanceFomP(q), "!=", p.Distance(q))
	var origin g.Point
	fmt.Println(distanceFomP(origin))

	scaleP := p.ScaleBy

	fmt.Println()
	scaleP(2)
	fmt.Println(p)
	scaleP(3)
	fmt.Println(p)
	scaleP(10)
	fmt.Println(p)
}

func f1() {
	var cp g.ColoredPoint
	cp.X = 100
	cp.Y = 200
	fmt.Println(cp)

	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	var p = g.ColoredPoint{g.Point{1, 1}, red}
	var q = g.ColoredPoint{g.Point{5, 4}, blue}

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
