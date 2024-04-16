package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
)

type IError struct {
	Field string `json:"field,omitempty"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value,omitempty"`
}

var Validator = validator.New()

func Validate(s any, c *fiber.Ctx) error {
	var errs []*IError

	err := json.Unmarshal(c.Body(), &s)

	if err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	err = Validator.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el IError
			pattern := regexp.MustCompile("(\\p{Lu}+\\P{Lu}*)")
			s = pattern.ReplaceAllString(err.Field(), "${1}_")
			s, _ = strings.CutSuffix(strings.ToLower(err.Field()), "_")
			el.Field = s.(string)
			el.Tag = err.Tag()
			el.Value = err.Param()
			errs = append(errs, &el)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"type":   "validator/list:validation",
			"errors": errs,
		})
	}

	c.Locals("requestValidation", s)
	return c.Next()
}
