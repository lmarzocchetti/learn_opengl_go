package pkg

import (
	"image"
	"image/draw"
	"log"
	"os"
)

func OpenImageRGBA(filename string) (*image.RGBA, int32, int32, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, 0, 0, err
	}

	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(reader)

	imageYCbCr, _, err := image.Decode(reader)
	if err != nil {
		return nil, 0, 0, err
	}

	width, height := imageYCbCr.Bounds().Size().X, imageYCbCr.Bounds().Size().Y
	imageRGBA := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(imageRGBA, imageRGBA.Bounds(), imageYCbCr, imageRGBA.Bounds().Min, draw.Src)

	return imageRGBA, int32(width), int32(height), nil
}
