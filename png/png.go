// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// The Image represents a structure for working with PNG images.
type Image struct {
	in  image.Image
	mid *image.RGBA64
	out *image.RGBA64

}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
func Load(filePath string) (*Image, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inImg, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	inBounds := inImg.Bounds()

	outImg := image.NewRGBA64(inBounds)

	midImg := image.NewRGBA64(inBounds)

	InitForMid(inImg, midImg)

	return &Image{inImg, midImg, outImg}, nil
}

func (img *Image)Swap(){
	img.mid = img.out
	img.out = image.NewRGBA64(img.in.Bounds())
}

func InitForMid(inImg image.Image, midImg *image.RGBA64) {
	inBounds := inImg.Bounds()
	for y := inBounds.Min.Y; y < inBounds.Max.Y; y++ {
		for x := inBounds.Min.X; x < inBounds.Max.X; x++ {
			r, g, b, a := inImg.At(x, y).RGBA()
			midImg.Set(x, y, color.RGBA64{clamp(float64(r)), clamp(float64(g)), clamp(float64(b)), uint16(a)})
		}
	}
}

func (img *Image) GetSubImage(x0 int, y0 int, x1 int, y1 int) image.Image{
	rect := image.Rect(x0,y0, x1, y1)
	return img.mid.SubImage(rect)
}

// Save saves the image to the given file
func (img *Image) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.out)
	if err != nil {
		return err
	}
	return nil
}

func (img *Image) Contains(x int, y int) bool {


	bounds := img.out.Bounds()

	if x < bounds.Min.X || x  >= bounds.Max.X{
		return false
	}
	if y < bounds.Min.Y || y  >= bounds.Max.Y{
		return false
	}
	return true
}

func (img *Image) GetSize() (int, int){
	bounds := img.in.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	return width, height
}



//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}
