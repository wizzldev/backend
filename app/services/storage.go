package services

import (
	"bytes"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
	"image/png"
	"io"
	"mime"
	"os"
)

type storage struct{}

var Storage storage

func (storage) WebPFromFormFile(file io.Reader, dest *os.File) error {
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	img = resize.Resize(512, 512, img, resize.Lanczos3)

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	return webp.Encode(dest, img, options)
}

func (storage) GetFileName(key, mimeType string) string {
	var t string
	types, err := mime.ExtensionsByType(mimeType)
	if err == nil {
		if len(types) > 0 {
			t = types[0]
		}
	}
	return key + t
}

func (storage) WebPStream(file io.Reader, size uint) (io.Reader, error) {
	img, err := webp.Decode(file, &decoder.Options{})
	if err != nil {
		return nil, err
	}

	if size >= 15 && size <= 1024 {
		img = resize.Resize(size, size, img, resize.Lanczos3)
	}

	buf := new(bytes.Buffer)
	err = webp.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
