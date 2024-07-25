package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/pkg/utils/role"
	"reflect"
	"regexp"
	"strings"
	"time"
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
	fmt.Println("validating")
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
			pattern := regexp.MustCompile(`([a-z0-9])([A-Z])`)
			p := pattern.ReplaceAllString(fieldError.Field(), "${1}_${2}")
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
	err := Validator.RegisterValidation("is_role", func(level validator.FieldLevel) bool {
		value := level.Field().String()
		_, err := role.New(strings.ToUpper(value))
		return err == nil
	})

	if err != nil {
		log.Fatal("Failed to register validator (is_role):", err)
	}

	err = Validator.RegisterValidation("invite_date", func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if !ok {
			return false
		}
		if date.IsZero() {
			return true
		}
		now := time.Now()
		return date.After(now.AddDate(0, 0, 1)) && date.Before(now.AddDate(0, 6, 0))
	})

	if err != nil {
		log.Fatal("Failed to register validator (invite_date):", err)
	}

	err = Validator.RegisterValidation("is_pointer", func(fl validator.FieldLevel) bool {
		return fl.Field().Kind() == reflect.Ptr
	})

	if err != nil {
		log.Fatal("Failed to register validator (is_pointer):", err)
	}

	err = Validator.RegisterValidation("is_emoji", func(fl validator.FieldLevel) bool {
		fmt.Println("validating is emoji")
		s := fl.Field().String()
		return IsEmoji(s)
	})

	if err != nil {
		log.Fatal("Failed to register validator (is_emoji):", err)
	}
}
