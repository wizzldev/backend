package services

import (
	"bytes"
	"errors"
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
	"image"
	"image/gif"
	"image/jpeg"
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

func (*Storage) WebPFromFormFile(file io.Reader, dest *os.File, contentType string) error {
	var (
		img image.Image
		err error
	)

	switch contentType {
	case "image/png":
		img, err = png.Decode(file)
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/gif":
		img, err = gif.Decode(file)
	case "image/webp":
		img, err = webp.Decode(file, &decoder.Options{})
	default:
		img, err = nil, errors.New("unsupported image type")
	}

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
	file := c.Locals(configs.LocalFileModel).(*models.File)
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
	if err != nil {
		return nil, err
	}
	defer file.Close()

	disc := s.NewDiscriminator()
	cType := fileH.Header.Get("Content-Type")
	path := s.getFileName(disc, "image/webp")
	dest, err := os.Create(filepath.Join(s.BasePath, path))
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	fileInfo, err := dest.Stat()
	if err != nil {
		return nil, err
	}

	err = s.WebPFromFormFile(file, dest, cType)
	if err != nil {
		return nil, err
	}

	fileModel := models.File{
		Path:          path,
		Name:          "",
		Type:          "avatar",
		Discriminator: disc,
		ContentType:   "image/webp",
		Size:          fileInfo.Size(),
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
	if err != nil {
		fmt.Println("failed to open new file", err)
		return nil, err
	}
	defer dest.Close()

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
		Size:          fileH.Size,
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

func (s *Storage) RemoveByDisc(disc string) error {
	var file models.File
	database.DB.Model(&models.File{}).Where("discriminator = ?", disc).First(&file)

	if file.ID < 1 {
		return errors.New("file not found")
	}

	err := os.Remove(filepath.Join(s.BasePath, file.Path))
	if err != nil {
		return err
	}

	database.DB.Delete(&file)

	return nil
}

func (s *Storage) OpenFile(path string) (*os.File, error) {
	return os.Open(filepath.Join(s.BasePath, path))
}
