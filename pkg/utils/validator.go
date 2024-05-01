package utils

import (
	"encoding/json"
	"errors"
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

func Validate[T any](c *fiber.Ctx) error {
	var s T
	var errs []*IError

	err := json.Unmarshal(c.Body(), &s)
	if err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	err = Validator.Struct(s)

	if err != nil {
		var validatorError validator.ValidationErrors
		if !errors.As(err, &validatorError) {
			return err
		}

		for _, fieldError := range validatorError {
			var el IError
			pattern := regexp.MustCompile("(\\p{Lu}+\\P{Lu}*)")
			p := pattern.ReplaceAllString(fieldError.Field(), "${1}_")
			p, _ = strings.CutSuffix(strings.ToLower(fieldError.Field()), "_")
			el.Field = p
			el.Tag = fieldError.Tag()
			el.Value = fieldError.Param()
			errs = append(errs, &el)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"type":   "error:validator",
			"errors": errs,
		})
	}

	c.Locals("requestValidation", &s)
	return c.Next()
}
