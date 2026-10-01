package bubble

import (
	"errors"
	"image"
	"image/draw"

	"forge.pi.home.arpa/govalds/bot/pkg/imaging"
	drawx "golang.org/x/image/draw"
)

var (
	ErrImageTooBig      = errors.New("image is too big, max size is 8192x8192")
	ErrMaskSizeMismatch = errors.New("alpha and overlay not same size")
	ErrMaskNil          = errors.New("both alpha and overlay are nil")
)

type position uint8

const (
	TopRight position = iota
	TopLeft
	BottomRight
	BottomLeft
)

type Mask interface {
	Overlay() *image.NRGBA
	Alpha() *image.Alpha
}

type mask struct {
	overlay *image.NRGBA
	alpha   *image.Alpha
}

func (m *mask) Overlay() *image.NRGBA { return m.overlay }
func (m *mask) Alpha() *image.Alpha   { return m.alpha }

func NewMask(overlay *image.NRGBA, alpha *image.Alpha) (Mask, error) {
	if overlay == nil && alpha == nil {
		return nil, ErrMaskNil
	}
	if overlay != nil && alpha != nil && overlay.Rect != alpha.Rect {
		return nil, ErrMaskSizeMismatch
	}
	return &mask{overlay, alpha}, nil
}

func Speechbubble(source *image.NRGBA, mask Mask, pos position, extend bool, maskHeight int) (*image.NRGBA, error) {
	overlay := mask.Overlay()
	alpha := mask.Alpha()
	bmask := overlay.Bounds()
	bsource := source.Bounds()

	flipH := pos == TopLeft || pos == BottomLeft
	flipV := pos == BottomLeft || pos == BottomRight

	if flipH || flipV {
		if overlay != nil {
			overlay = imaging.FlipNRGBA(overlay, flipV, flipH)
		}
		if alpha != nil {
			alpha = imaging.FlipAlpha(alpha, flipV, flipH)
		}
	}

	var height int

	if maskHeight != 100 {
		height = source.Bounds().Dy() / 100 * maskHeight
	} else {
		height = source.Bounds().Dy()
	}

	newRect := image.Rect(bsource.Min.X, bsource.Min.Y, bsource.Dx(), height)

	if alpha != nil {
		alphaNew := image.NewAlpha(newRect)
		drawx.CatmullRom.Scale(alphaNew, newRect, alpha, bmask, draw.Over, nil)
		alpha = alphaNew
	}

	if overlay != nil {
		overlayNew := image.NewNRGBA(newRect)
		drawx.CatmullRom.Scale(overlayNew, newRect, overlay, bmask, draw.Over, nil)
		overlay = overlayNew
	}

	if overlay != nil {
		bmask = overlay.Bounds()
	} else {
		bmask = alpha.Bounds()
	}

	if extend {
		if pos == TopLeft || pos == TopRight {
			source = fillTop(source, bmask)
		} else {
			source = fillBottom(source, bmask)
		}
	}

	var bottom bool
	if pos == BottomLeft || pos == BottomRight {
		bottom = true
	}

	return drawBubble(source, overlay, alpha, bottom)
}
