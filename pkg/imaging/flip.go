package imaging

import "image"

func FlipAlpha(source *image.Alpha, vertical, horizontal bool) *image.Alpha {
	if vertical == false && horizontal == false {
		return source
	}

	b := source.Bounds()
	w, h := b.Dx(), b.Dy()

	flipped := image.NewAlpha(b)

	for y := range h {
		for x := range w {
			destX := x
			destY := y

			if vertical {
				destY = h - 1 - y
			}
			if horizontal {
				destX = w - 1 - x
			}

			flipped.SetAlpha(destX, destY, source.AlphaAt(x, y))
		}
	}
	return flipped
}

func FlipNRGBA(source *image.NRGBA, vertical, horizontal bool) *image.NRGBA {
	if vertical == false && horizontal == false {
		return source
	}

	b := source.Bounds()
	w, h := b.Dx(), b.Dy()

	flipped := image.NewNRGBA(b)

	for y := range h {
		for x := range w {
			destX := x
			destY := y

			if vertical {
				destY = h - 1 - y
			}
			if horizontal {
				destX = w - 1 - x
			}

			flipped.SetNRGBA(destX, destY, source.NRGBAAt(x, y))
		}
	}
	return flipped
}
