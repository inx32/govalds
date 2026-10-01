package imaging

import (
	"image"
	"image/draw"
)

func ToNRGBA(img image.Image) *image.NRGBA {
	b := img.Bounds()
	output := image.NewNRGBA(b)
	draw.Draw(output, b, img, b.Min, draw.Src)
	return output
}

func ToAlpha(img image.Image) *image.Alpha {
	b := img.Bounds()
	output := image.NewAlpha(b)
	draw.Draw(output, b, img, b.Min, draw.Src)
	return output
}
