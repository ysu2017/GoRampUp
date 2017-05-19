package imageUtils

import (
	"image"
	"image/color"
	"yi/point"
)

func ImageToGrey(img image.Image) *image.Gray {
	newRec := img.Bounds()
	greyImg := image.NewGray(newRec)
	for i := newRec.Min.X; i <= newRec.Max.X; i++ {
		for j := newRec.Min.Y; j <= newRec.Max.Y; j++ {
			greyImg.Set(i, j, img.At(i, j))
		}
	}
	return greyImg
}

func ImageToRGBA64(img image.Image) *image.RGBA64 {
	newRec := img.Bounds()
	rgbaImg := image.NewRGBA64(newRec)
	for i := newRec.Min.X; i <= newRec.Max.X; i++ {
		for j := newRec.Min.Y; j <= newRec.Max.Y; j++ {
			rgbaImg.Set(i, j, img.At(i, j))
		}
	}
	return rgbaImg
}

func MarkPointOnRGBA64(p point.Point, img *image.RGBA64) {
	red := color.RGBA{
		R: 1<<8 - 1,
		G: 0,
		B: 0,
		A: 1<<8 - 1,
	}
	for i := p.X - 10; i <= p.X+10; i++ {
		img.Set(i, p.Y, red)
		img.Set(i, p.Y+1, red)
	}
	for j := p.Y - 10; j <= p.Y+10; j++ {
		img.Set(p.X, j, red)
		img.Set(p.X+1, j, red)
	}
}
