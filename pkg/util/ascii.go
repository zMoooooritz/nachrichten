package util

import (
	"image"
	"math"

	"github.com/aybabtme/rgbterm"
	"github.com/nfnt/resize"

	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

const (
	ASCII_ASPECT_RATION = 4.0 / 1.0
)

func getChar(val float64) byte {
	asciiMatrix := []byte(" `^\":;Il!iXYUJCLQ0OZ#MW8B@$")

	divisor := float64(256) / float64(len(asciiMatrix))

	return asciiMatrix[int(math.Floor(float64(val)/divisor))]
}

func getLuminosityPt(x int, y int, img image.Image) float64 {
	r, g, b, _ := img.At(x, y).RGBA()
	return (0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
}

func applyColor(r uint8, g uint8, b uint8, character byte) string {
	return rgbterm.FgString(string([]byte{character}), r, g, b)
}

func ImageToAscii(img image.Image, width uint, height uint, colored bool) [][]string {
	w := width
	h := uint(float64(width) / ASCII_ASPECT_RATION)

	if h > height {
		w = uint(height * uint(ASCII_ASPECT_RATION))
		h = height
	}

	img = resize.Resize(w, h, img, resize.Lanczos3)
	bounds := img.Bounds()

	asciiImage := make([][]string, img.Bounds().Max.Y)
	for i := range asciiImage {
		asciiImage[i] = make([]string, img.Bounds().Max.X)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			lum := getLuminosityPt(x, y, img)
			character := getChar(lum)

			r, g, b, _ := color.NRGBAModel.Convert(img.At(x, y)).RGBA()

			if colored {
				asciiImage[y][x] = applyColor(uint8(r), uint8(g), uint8(b), character)
			} else {
				asciiImage[y][x] = string(character)
			}
		}
	}
	return asciiImage
}
