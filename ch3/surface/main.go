package main

import (
	"fmt"
	"math"
	"net/http"
)

const (
	width, height = 600, 320            // canvas size in px
	cells         = 100                 // number of gride cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         //angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

// const l = 16711680 - 255

// func xx() [l]int {
// 	var c [l]int
// 	for i := 0x0000ff; i < 0xff0000; i++ {
// 		c[i-0x0000ff] = i
// 	}
// 	// for i := 255; i < 1000; i++ {
// 	// 	fmt.Printf("\n%#x", i)
// 	// }
// 	return c
// }

func main() {
	http.HandleFunc("/h", handler)

	http.ListenAndServe(":3030", nil)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "image/svg+xml")

	svg := surface()
	fmt.Fprint(rw, svg)
	return
}

func surface() string {
	svg := fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" "+
		"style=\"stroke: grey; fill: white; stroke-width: 0.7\" "+
		"width=\"%d\" height=\"%d\">", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ay) || math.IsNaN(by) || math.IsNaN(cy) {
				continue
			}
			svg += fmt.Sprintf("<polygon points=\"%g,%g %g,%g %g,%g %g,%g\" style=\"fill:white;\" />\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}

	svg += fmt.Sprint("</svg>")

	return svg
}

func corner(i, j int) (float64, float64) {
	// find point (x, y) at corner of cells (i, j)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// compute surfave height z.
	z := f(x, y)

	// project (x,y,z) isomerically onto 2-D SVG canvas (sx, sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}
