package main

import (
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"yi/harris"
	"yi/point"
)

func markPoint(p point.Point, img *image.RGBA64) {
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

func main() {
	log.SetOutput(os.Stdout)
	output := flag.String("o", "new.jpg", "path to output image")
	flag.Parse()

	input := flag.Arg(0)
	if input == "" {
		log.Fatal("please provide a valid input image")
	}

	file, err := os.Open(input)
	if err != nil {
		log.Fatal("error while opening file")
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("error decoding image")
	}

	newRec := img.Bounds()
	newImg := image.NewRGBA64(img.Bounds())
	greyImg := image.NewGray(img.Bounds())

	for i := newRec.Min.X; i <= newRec.Max.X; i++ {
		for j := newRec.Min.Y; j <= newRec.Max.Y; j++ {
			newImg.Set(i, j, img.At(i, j))
			greyImg.Set(i, j, img.At(i, j))
		}
	}
	corners, err := harris.HarrisCornerDetector(greyImg, 1000000)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range corners {
		markPoint(p, newImg)
	}

	newFile, err := os.Create(*output)
	if err != nil {
		log.Fatal("Error creating new file")
	}
	defer newFile.Close()
	jpeg.Encode(newFile, newImg, nil)
}
