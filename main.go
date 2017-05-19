package main

import (
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"yi/harris"
	"yi/imageUtils"
)

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

	greyImg := imageUtils.ImageToGrey(img)
	rgba64Img := imageUtils.ImageToRGBA64(img)

	corners, err := harris.HarrisCornerDetector(greyImg, 1000000)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range corners {
		imageUtils.MarkPointOnRGBA64(p, rgba64Img)
	}

	newFile, err := os.Create(*output)
	if err != nil {
		log.Fatal("Error creating new file")
	}
	defer newFile.Close()
	jpeg.Encode(newFile, rgba64Img, nil)
}
