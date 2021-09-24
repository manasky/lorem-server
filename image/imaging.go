package image

import (
	"github.com/disintegration/imaging"
	"image"
)

type Imaging struct {}

func (i *Imaging) CropCenter(img image.Image, width, height int) *image.Image {
	img = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	return &img
}
