package services

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/utils"
	"image/png"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Storage struct {
	BasePath string
}

func NewStorage() (*Storage, error) {
	base, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	base = filepath.Join(base, "./storage")
	s := &Storage{
		BasePath: base,
	}
	return s, nil
}

func (*Storage) WebPFromFormFile(file io.Reader, dest *os.File) error {
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

func (*Storage) getFileName(key, mimeType string) string {
	var t string
	types, err := mime.ExtensionsByType(mimeType)
	if err == nil {
		if len(types) > 0 {
			t = types[0]
		}
	}
	return key + t
}

func (*Storage) WebPStream(file io.Reader, size uint) (io.Reader, error) {
	img, err := webp.Decode(file, &decoder.Options{})
	if err != nil {
		return nil, err
	}

	if size >= 1 && size <= 1024 {
		img = resize.Resize(size, size, img, resize.Lanczos3)
	}

	buf := new(bytes.Buffer)
	err = webp.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func (s *Storage) LocalFile(c *fiber.Ctx) *models.File {
	file := c.Locals(utils.LocalFileModel).(*models.File)
	file.Path = filepath.Join(s.BasePath, file.Path)
	return file
}

func (*Storage) SaveWebP(source io.Reader, dest *os.File) error {
	img, err := png.Decode(source)
	if err != nil {
		return err
	}

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)

	if err != nil {
		return err
	}

	return webp.Encode(dest, img, options)
}

func (*Storage) NewDiscriminator() string {
	rand := utils.NewRandom()
	key := rand.String(25)
	for {
		var count int64
		database.DB.Model(&models.File{}).Where("discriminator = ?", key).Count(&count)
		if count == 0 {
			break
		}
		key = rand.String(35)
	}
	return key
}

func (s *Storage) StoreAvatar(fileH *multipart.FileHeader) (*models.File, error) {
	file, err := fileH.Open()
	defer file.Close()
	if err != nil {
		return nil, err
	}

	disc := s.NewDiscriminator()
	path := s.getFileName(disc, fileH.Header.Get("Content-Type"))
	dest, err := os.Create(filepath.Join(s.BasePath, path))
	defer dest.Close()
	if err != nil {
		return nil, err
	}

	err = s.WebPFromFormFile(file, dest)
	if err != nil {
		return nil, err
	}

	fileModel := models.File{
		Path:          path,
		Name:          fileH.Filename,
		Type:          "avatar",
		Discriminator: disc,
		ContentType:   fileH.Header.Get("Content-Type"),
	}
	err = database.DB.Create(&fileModel).Error

	if err != nil {
		return nil, err
	}

	return &fileModel, nil
}

func (s *Storage) Store(fileH *multipart.FileHeader, token ...string) (*models.File, error) {
	if fileH.Size > configs.Env.MaxFileSize {
		return nil, fiber.NewError(fiber.StatusRequestEntityTooLarge, "file too large")
	}

	file, err := fileH.Open()
	if err != nil {
		fmt.Println("file header open error", err)
		return nil, err
	}

	disc := s.NewDiscriminator()
	path := s.getFileName(disc, fileH.Header.Get("Content-Type"))

	dest, err := os.Create(filepath.Join(s.BasePath, path))
	defer dest.Close()
	if err != nil {
		fmt.Println("failed to open new file", err)
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("failed to read file", err)
		return nil, err
	}

	_, err = dest.Write(data)
	if err != nil {
		fmt.Println("failed to write file", err)
		return nil, err
	}

	fileModel := models.File{
		Path:          path,
		Name:          fileH.Filename,
		Discriminator: disc,
		Type:          "file",
		ContentType:   fileH.Header.Get("Content-Type"),
	}

	if len(token) > 0 {
		t := &token[0]
		fileModel.AccessToken = t
	}

	err = database.DB.Create(&fileModel).Error
	if err != nil {
		return nil, err
	}

	return &fileModel, nil
}
