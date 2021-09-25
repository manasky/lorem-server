package image

import (
	"image"
	"os"
)

type Processor interface {
	CropCenter(img image.Image, width, height int) *image.Image
}

func Decode(name string) (*image.Image, error) {
	imagePath, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer imagePath.Close()

	img, _, err := image.Decode(imagePath)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
