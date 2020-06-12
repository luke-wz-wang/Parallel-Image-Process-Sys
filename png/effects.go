// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

var sharp = [9]float64{
	0, -1, 0,
	-1, 5, -1,
	0, -1, 0}

var edge = [9]float64{
	-1, -1, -1,
	-1, 8, -1,
	-1, -1, -1}

var blur = [9]float64{
	1. / 9, 1. / 9, 1. / 9,
	1. / 9, 1. / 9, 1. / 9,
	1. / 9, 1. / 9, 1. / 9}

// apply a certain effect to a rect area [(0,y0), (width, y1)] of the img
func (img *Image) AddEffect(effect string,  y0 int, y1 int){

	switch effect{
	case "S":
		img.FilterRect(sharp, y0, y1)
	case "E":
		img.FilterRect(edge, y0, y1)
	case "B":
		img.FilterRect(blur, y0, y1)
	case "G":
		img.GrayscaleRect(y0, y1)
	default:
		img.GrayscaleRect(y0, y1)
	}
}

// apply a certain effect to a rect area [(0,y0), (width, y1)] of the img
func (img *Image) AddEffectRect(effect string, y0 int, y1 int, ch chan bool){
	img.AddEffect(effect, y0, y1)
	ch <- true
}

// apply filter k to a rect area [(0,y0), (width, y1)] of the img
func (img *Image) FilterRect(k [9]float64, y0 int, y1 int){
	bounds := img.out.Bounds()
	for y := y0; y < y1; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sum_r float64
			var sum_g float64
			var sum_b float64
			var i int
			_, _, _, a := img.mid.At(x, y).RGBA()
			for yo := y - 1; yo <= y+1; yo++ {
				for xo := x - 1; xo <= x+1; xo++ {
					if img.Contains(xo, yo) {
						r, g, b, _ := img.mid.At(xo, yo).RGBA()
						sum_r += k[i] * float64(r)
						sum_g += k[i] * float64(g)
						sum_b += k[i] * float64(b)
					}
					i++
				}
			}
			img.out.Set(x, y, color.RGBA64{clamp(sum_r), clamp(sum_g), clamp(sum_b), uint16(a)})
		}
	}
}

// apply greyscale effect to a rect area [(0,y0), (width, y1)] of the img
func (img *Image) GrayscaleRect(y0 int, y1 int) {
	bounds := img.out.Bounds()
	for y := y0; y < y1; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.mid.At(x, y).RGBA()
			greyC := clamp(float64(r+g+b) / 3)
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}






