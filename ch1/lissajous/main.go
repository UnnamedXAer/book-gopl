// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
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

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
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
		for t := 0.0; t < cycles*2*math.Pi; t += res {
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
