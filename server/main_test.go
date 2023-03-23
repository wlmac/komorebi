package server

import (
	"embed"
	"image"
	_ "image/jpeg"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
)

const (
	width  = 1920
	height = 0 // preserve aspect ratio
)

// TODO: test imagemagick- and libvips- based library

//go:embed testdata/images
var images embed.FS

var paths = getPaths()

func getPaths() []string {
	matches, err := fs.Glob(images, "testdata/images/*")
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	for _, match := range matches {
		if strings.HasSuffix(match, ".attribution.txt") || strings.Contains(match, "README.md") {
			continue
		}
		res = append(res, match)
	}
	return res
}

func disintegration(b *testing.B, sourcePath string, width, height int, filter imaging.ResampleFilter) {
	f, err := images.Open(sourcePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	img, _, err := image.Decode(f)
	if err != nil {
		b.Fatalf("%s: %s", sourcePath, err)
	}
	resized := imaging.Resize(img, width, height, filter)
	encoder, err := encoder.NewEncoder(resized, webpEncodeOptions)
	if err != nil {
		panic(err)
	}
	err = encoder.Encode(io.Discard)
	if err != nil {
		panic(err)
	}
}

func BenchmarkDisintegrationLanczos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			disintegration(b, path, width, height, imaging.Lanczos)
		}
	}
}

func BenchmarkDisintegrationBox(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			disintegration(b, path, width, height, imaging.Box)
		}
	}
}

func BenchmarkDisintegrationLinear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			disintegration(b, path, width, height, imaging.Linear)
		}
	}
}
