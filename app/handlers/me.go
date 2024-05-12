package handlers

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizzldev/chat/database"
	"image"
	"image/png"
	"os"
	"strings"
)

type me struct{}

var Me me

func (me) Hello(c *fiber.Ctx) error {
	user := authUser(c)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Hello %s", user.FirstName),
		"user":    user,
	})
}

func (me) UploadProfileImage(c *fiber.Ctx) error {
	h, err := c.FormFile("image")

	if err != nil {
		return err
	}

	file, err := h.Open()

	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	img = imaging.Resize(img, 300, 300, imaging.Box)

	user := authUser(c)
	uniqueId := uuid.New()
	imageName := fmt.Sprintf("u%v-i%s.png", user.ID, strings.Replace(uniqueId.String(), "-", "", -1))
	out, err := os.Create(fmt.Sprintf("./storage/image/%s", imageName))

	if err != nil {
		return err
	}

	err = png.Encode(out, img)

	if err != nil {
		return err
	}

	if user.ImageURL != "" {
		splitImage := strings.Split(user.ImageURL, "/")
		imgName := splitImage[len(splitImage)-1]
		imgPath := fmt.Sprintf("./storage/image/%s", imgName)
		_ = os.Remove(imgPath)
	}

	user.ImageURL = fmt.Sprintf("{cdn}/static/image/%s", imageName)
	database.DB.Save(user)

	return c.JSON(user)
}
