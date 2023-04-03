package server

import (
	"image"
	"io"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

var webpDecodeOptions = &decoder.Options{}

var webpEncodeOptions = &encoder.Options{}

func init() {
	image.RegisterFormat("webp", "RIFF????WEBP", webpDecode, webpDecodeConfig)
}

func webpDecode(r io.Reader) (image.Image, error) {
	return webp.Decode(r, webpDecodeOptions)
}

func webpDecodeConfig(r io.Reader) (image.Config, error) {
	img, err := webp.Decode(r, webpDecodeOptions)
	if err != nil {
		return image.Config{}, err
	}
	return image.Config{
		ColorModel: img.ColorModel(),
		Width:      img.Bounds().Dx(),
		Height:     img.Bounds().Dy(),
	}, nil
}
