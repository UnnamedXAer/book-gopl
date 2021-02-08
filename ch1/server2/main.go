// Lock resources
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var palette = []color.Color{
	color.White,
	color.RGBA{
		R: 210,
		G: 111,
		B: 240,
		A: 1,
	},
	color.RGBA{
		R: 100,
		G: 111,
		B: 240,
		A: 1,
	},
	color.RGBA{
		R: 99,
		G: 212,
		B: 114,
		A: 1,
	},
	color.RGBA{
		R: 238,
		G: 217,
		B: 83,
		A: 1,
	}}

const (
	bgColorIdx = 0
	fgColorIdx = 1
)

var (
	count int64
	mu    sync.Mutex
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", countHandler)
	http.HandleFunc("/l", lissajousHandler)
	http.HandleFunc("/exit", exitHandler)

	fmt.Println(http.ListenAndServe(":3030", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}

	// fmt.Fprintf(w, "RequestURI: %q\nURL: %q\n", r.RequestURI, r.URL.Path)
}
func countHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "RequestURI: %q\nURL: %q\n\n\tcount: %d", r.RequestURI, r.URL.Path, count)
	mu.Unlock()
}
func exitHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Printf("RequestURI: %q\nURL: %q\n\n\tcount: %d\n\n\n\nShutting down server...", r.RequestURI, r.URL.Path, count)
	mu.Unlock()

	os.Exit(0)
}
func lissajousHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(404)
		fmt.Fprint(w, err)
		return
	}

	cycles := r.FormValue("cycles")
	cNum, err := strconv.ParseFloat(cycles, 64)
	if err != nil {
		fmt.Fprintf(w, "RequestURI: %q\nURL: %q\n\n\nwrong value of the query param 'cycles' = %q", r.RequestURI, r.URL.Path, cycles)
		return
	}
	if cNum > 500 || cNum < 1 {
		fmt.Fprintf(w, "RequestURI: %q\nURL: %q\n\n\nwrong value of the query param 'cycles' = %s, allowed range [1:500]", r.RequestURI, r.URL.Path, cycles)
		return
	}

	lissajous(w, cNum)
}

func lissajous(out io.Writer, cycles float64) {
	const (
		// cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 //angular resolution
		size    = 100   // image canvas covers [-size...+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	colorIdx := 0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < (cycles)*2*math.Pi; t += res {
			if int(t*1000)%1000 == 0 {
				colorIdx++
				if colorIdx > 4 {
					colorIdx = 1
				}
			}
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(colorIdx))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	err := gif.EncodeAll(out, &anim)
	if err != nil {
		fmt.Println(err)
	}
}
