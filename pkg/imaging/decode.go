package imaging

import (
	"bytes"
	"errors"
	"image"
	"io"
)

const imageMaxBytes = 8 * 1024 * 1024
const imageMaxPixels = 2048 * 2048
const imageMaxDimension = 4096

var ErrFileTooBig = errors.New("file too big")
var ErrImageTooBig = errors.New("image too big")
var ErrImageInvalid = errors.New("invalid image")

func Decode(source io.Reader) (image.Image, error) {
	data, err := io.ReadAll(io.LimitReader(source, imageMaxBytes+1))
	if err != nil {
		return nil, err
	}

	if len(data) > imageMaxBytes {
		return nil, ErrFileTooBig
	}

	r := bytes.NewReader(data)

	cfg, format, err := image.DecodeConfig(r)
	if err != nil {
		return nil, err
	}

	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, ErrImageInvalid
	}

	w, h := int64(cfg.Width), int64(cfg.Height)
	if w > imageMaxPixels/h || w > imageMaxDimension || h > imageMaxDimension {
		return nil, ErrImageTooBig
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	img, decodedFormat, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	if decodedFormat != format {
		return nil, ErrImageInvalid
	}

	return img, nil
}
