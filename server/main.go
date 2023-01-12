package server

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
)

var SupportedFormats = map[string]struct{}{
	"webp": struct{}{},
	"jpeg": struct{}{},
	"png":  struct{}{},
}

type Config struct {
	SourcePath string
	// SourcePath specifies path of dir where source media are stored under.

	CachePath string
	// CachePath specifies path of directory to store cached media under.
}

type server struct {
	c Config
}

func New(config Config) (http.Handler, error) {
	if config.CachePath == "" {
		return nil, errors.New("CachePath must not be blank")
	}
	return &server{
		c: config,
	}, nil
}

type editConfig struct {
	Width, Height int
}

// getMedia gets an edited image from using sourcePath and cachePath.
// The caller must close the returned io.ReadCloser if it is not nil.
func (s *server) getMedia(ec editConfig, w io.Writer, sourcePath, cachePath string) (err error) {
	cached, err := os.Open(cachePath)
	if errors.Is(err, fs.ErrNotExist) {
		// resize and serve
		var cache *os.File
		cache, err = os.Create(cachePath)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		defer func() {
			closeErr := cache.Close()
			if closeErr != nil {
				err = closeErr
			}
		}()
		w2 := io.MultiWriter(w, cache)
		defer func() {
			if err := recover(); err != nil {
				rmErr := os.Remove(cachePath)
				if rmErr != nil {
					log.Printf("remove cache: %S", rmErr)
				}
			}
		}()
		err = s.editMedia(ec, w2, sourcePath)
		if err != nil {
			rmErr := os.Remove(cachePath)
			if rmErr != nil {
				err = rmErr
			}
		}
		return
	} else if err != nil {
		return
	}
	_, err = io.Copy(w, cached)
	return
}

func (s *server) editMedia(ec editConfig, w io.Writer, sourcePath string) (err error) {
	f, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()
	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	resized := imaging.Resize(img, ec.Width, ec.Height, imaging.Lanczos)
	encoder, err := encoder.NewEncoder(resized, webpEncodeOptions)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	err = encoder.Encode(w)
	return
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	rid := ""
	rid2 := make([]byte, 16)
	_, err := rand.Read(rid2)
	if err != nil {
		rid = "error"
	} else {
		rid = base64.StdEncoding.EncodeToString(rid2)
	}
	w.Header().Add("Server", "Komorebi/0")
	w.Header().Add("X-Request-ID", rid)
	if err != nil {
		log.Printf("request %s: error %s", rid, err)
	}

	q := r.URL.Query()
	width := q.Get("w")
	height := q.Get("h")
	if width != "" && height != "" {
		w.WriteHeader(422)
		fmt.Fprint(w, "cannot specify both width and height")
		return
	}
	var width2, height2 int64
	if width == "" {
		width2 = 0
	} else {
		width2, err = strconv.ParseInt(width, 10, 32)
		if err != nil {
			w.WriteHeader(422)
			fmt.Fprint(w, "cannot parse width")
			return
		}
	}
	if height == "" {
		height2 = 0
	} else {
		height2, err = strconv.ParseInt(height, 10, 32)
		if err != nil {
			w.WriteHeader(422)
			fmt.Fprint(w, "cannot parse height")
			return
		}
	}
	formats := q["fmt"]
	if len(formats) == 0 {
		err = errors.New("must specify formats")
		w.WriteHeader(422)
		fmt.Fprint(w, "must specify formats")
		return
	}
	if len(formats) > len(SupportedFormats) {
		err = errors.New("too many formats")
		w.WriteHeader(422)
		fmt.Fprint(w, "too many formats")
		return
	}
	for _, format := range formats {
		_, ok := SupportedFormats[format]
		if ok {
			goto SupportedFormatFound
		}
	}
	err = fmt.Errorf("only unsupported formats (%s) specified", formats)
	w.WriteHeader(422)
	fmt.Fprint(w, "only unsupported formats specified")
	return
SupportedFormatFound:
	var sourcePath string
	sourcePath, err = url.JoinPath("/", r.URL.Path)
	if err != nil {
		err = fmt.Errorf("sourcePath: %w", err)
		w.WriteHeader(422)
		fmt.Fprint(w, "url invalid")
		return
	}
	sourcePath = filepath.Join(s.c.SourcePath, sourcePath)

	cachePath := filepath.Join(s.c.CachePath, base64.URLEncoding.EncodeToString([]byte(sourcePath)))
	w.WriteHeader(200)
	// TODO: Content-Type
	err = s.getMedia(editConfig{
		Width:  int(width2),
		Height: int(height2),
	}, w, sourcePath, cachePath)
	if err != nil {
		err = fmt.Errorf("getMedia: %w", err)
		fmt.Fprint(w, "getting media failed")
		return
	}
}
