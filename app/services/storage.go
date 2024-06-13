package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image/png"
	"io"
	"mime"
	"os"
)

type storage struct{}

var Storage storage

func (storage) WebPFromFormFile(c *fiber.Ctx, file io.Reader, dest *os.File) error {
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)

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
