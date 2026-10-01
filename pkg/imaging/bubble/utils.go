package bubble

import (
	"image"
	"image/draw"
)

func fillTop(source *image.NRGBA, bmask image.Rectangle) *image.NRGBA {
	bsource := source.Bounds()

	maskWidth := bsource.Dx()
	maskHeight := bmask.Dy() * maskWidth / bmask.Dx()

	output := image.NewNRGBA(image.Rectangle{
		Min: bsource.Min,
		Max: image.Pt(bsource.Max.X, bsource.Max.Y+maskHeight),
	})

	sourceRect := image.Rectangle{
		Min: image.Pt(0, maskHeight),
		Max: image.Pt(bsource.Dx(), maskHeight+bsource.Dy()),
	}
	draw.Draw(output, sourceRect, source, bsource.Min, draw.Src)

	for x := range bsource.Dx() {
		pixel := source.At(x, bsource.Min.Y)

		fillRect := image.Rectangle{
			Min: image.Pt(x, 0),
			Max: image.Pt(x+1, maskHeight),
		}
		fillSrc := &image.Uniform{pixel}
		draw.Draw(output, fillRect, fillSrc, image.Pt(0, 0), draw.Src)
	}

	return output
}

func fillBottom(source *image.NRGBA, bmask image.Rectangle) *image.NRGBA {
	bsource := source.Bounds()

	maskWidth := bsource.Dx()
	maskHeight := bmask.Dy() * maskWidth / bmask.Dx()

	output := image.NewNRGBA(image.Rectangle{
		Min: bsource.Min,
		Max: image.Pt(bsource.Max.X, bsource.Max.Y+maskHeight),
	})

	draw.Draw(output, bsource, source, bsource.Min, draw.Src)

	for x := 0; x < bsource.Dx(); x++ {
		pixel := source.At(x, bsource.Max.Y-1)

		fillRect := image.Rectangle{
			Min: image.Pt(x, bsource.Max.Y),
			Max: image.Pt(x+1, bsource.Max.Y+maskHeight),
		}
		fillSrc := &image.Uniform{pixel}
		draw.Draw(output, fillRect, fillSrc, image.Pt(0, 0), draw.Src)
	}

	return output
}

func drawBubble(source, overlay *image.NRGBA, alpha *image.Alpha, bottom bool) (*image.NRGBA, error) {
	if overlay != nil && alpha != nil && overlay.Rect != alpha.Rect {
		return nil, ErrMaskSizeMismatch
	}

	bsource := source.Bounds()
	var bmask image.Rectangle

	if overlay != nil {
		bmask = overlay.Bounds()
	} else {
		bmask = alpha.Bounds()
	}

	output := image.NewNRGBA(bsource)
	draw.Draw(output, bsource, source, bsource.Min, draw.Src)

	offsetY := 0
	if bottom {
		offsetY = bsource.Dy() - bmask.Dy()
	}

	if alpha != nil {
		for x := range bmask.Dx() {
			for y := range bmask.Dy() {
				destY := y + offsetY

				alphaPixel := alpha.AlphaAt(x, y)
				sourcePixel := source.NRGBAAt(x, destY)

				if sourcePixel.A != alphaPixel.A {
					sourcePixel.A = alphaPixel.A
					output.SetNRGBA(x, destY, sourcePixel)
				}
			}
		}
	}

	if overlay != nil {
		destRect := image.Rect(0, offsetY, bmask.Dx(), offsetY+bmask.Dy())
		draw.Draw(output, destRect, overlay, image.Point{}, draw.Over)
	}

	return output, nil
}
