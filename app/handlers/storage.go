package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/utils"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

type storage struct {
	BasePath string
}

var Storage = &storage{}

func (s *storage) Init() error {
	base, err := os.Getwd()
	if err != nil {
		return err
	}
	base = filepath.Join(base, "./storage")
	s.BasePath = base
	return nil
}

func (s *storage) Get(c *fiber.Ctx) error {
	return c.SendFile(s.localFile(c).Path)
}

func (s *storage) GetAvatar(c *fiber.Ctx) error {
	fileModel := s.localFile(c)
	if fileModel.Type != "avatar" {
		return fiber.ErrBadRequest
	}

	size := c.QueryInt("s", 0)

	file, err := os.Open(fileModel.Path)
	if err != nil {
		return err
	}

	stream, err := services.Storage.WebPStream(file, uint(size))
	if err != nil {
		return err
	}

	c.Response().Header.Set("Content-Type", "image/webp")
	c.Response().Header.Set("Cache-Control", "max-age=3600")
	return c.SendStream(stream)
}

func (s *storage) StoreAvatar(c *fiber.Ctx) error {
	fileH, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	file, err := fileH.Open()
	defer file.Close()
	if err != nil {
		return err
	}

	disc := s.newDiscriminator()
	path := services.Storage.GetFileName(disc, fileH.Header.Get("Content-Type"))
	dest, err := os.Open(filepath.Join(s.BasePath, path))
	defer dest.Close()
	if err != nil {
		return err
	}

	err = services.Storage.WebPFromFormFile(file, dest)
	if err != nil {
		return err
	}

	err = database.DB.Create(&models.File{
		Path:          path,
		Name:          fileH.Filename,
		Type:          "avatar",
		Discriminator: disc,
		ContentType:   fileH.Header.Get("Content-Type"),
	}).Error
	if err != nil {
		return err
	}

	user := authUser(c)
	user.ImageURL = fmt.Sprintf("/storage/avatars/%s", path)

	database.DB.Save(user)

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (s *storage) localFile(c *fiber.Ctx) *models.File {
	file := c.Locals(utils.LocalFileModel).(*models.File)
	fmt.Println(s.BasePath)
	file.Path = filepath.Join(s.BasePath, file.Path)
	fmt.Println("filepath", file.Path)
	return file
}

func (*storage) saveWebP(source io.Reader, dest *os.File) error {
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

func (*storage) newDiscriminator() string {
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
