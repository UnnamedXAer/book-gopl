package geometry

import (
	"image/color"
	"math"
)

// Point represents point in two dimensional space
type Point struct {
	X float64
	Y float64
}

// Distance returns distance between p and q points
// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Distance returns distance between p and q points
// same functionality, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

// Path is a journey connecting the points with straight lines.
type Path []Point

func (path Path) Distance() float64 {
	count := len(path)
	if count == 0 {
		return 0.0
	}
	if count == 1 {
		return 0.0
	}

	var d float64
	for i := 0; i < count-1; i++ {
		d += path[i].Distance(path[i+1])
	}
	return d
}

func (path Path) TranslateBy(offset Point, add bool) {
	var op func(p, q Point) Point
	if add {
		op = Point.Add
	} else {
		op = Point.Sub
	}

	for i := range path {
		// call either path[i].Add(offser) or path[i].Sub(offset)
		path[i] = op(path[i], offset)
	}
}

// ColoredPoint represents point in two dimensional space
// with additional color value
type ColoredPoint struct {
	Point
	Color color.RGBA
}

// func main() {
// 	fmt.Println(Path{
// 		{1, 1},
// 		{5, 1},
// 		{5, 4},
// 		{1, 1},
// 	}.Distance())

// 	r := &Point{1, 23}
// 	r.ScaleBy(2)
// 	fmt.Println(*r)
// }
