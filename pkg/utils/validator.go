package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/utils/role"
	"regexp"
	"strings"
)

type IError struct {
	Field string `json:"field,omitempty"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value,omitempty"`
}

var (
	Validator    = validator.New()
	IsRegistered = false
)

func Validate[T any](c *fiber.Ctx) error {
	if !IsRegistered {
		RegisterCustomValidations()
	}
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
			pattern := regexp.MustCompile(`(\\p{Lu}+\\P{Lu}*)`)
			p := pattern.ReplaceAllString(fieldError.Field(), "${1}_")
			p = strings.TrimRight(p, "_")
			el.Field = strings.ToLower(p)
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

func RegisterCustomValidations() {
	IsRegistered = true
	_ = Validator.RegisterValidation("is_role", func(level validator.FieldLevel) bool {
		value := level.Field().String()
		_, err := role.New(strings.ToUpper(value))
		return err == nil
	})
}
