package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/tests"
	"net/http/httptest"
	"testing"
)

func Test_Register(t *testing.T) {
	app := tests.NewApp("../..")

	// Connect to a test db
	database.MustConnectTestDB()

	app.Post("/", requests.Use[requests.Register](), Auth.Register)

	data, _ := json.Marshal(fiber.Map{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john.doe@example.com",
		"password":   "secret1234",
	})

	req := httptest.NewRequest("POST", "/", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	if err != nil {
		t.Fatal("HTTP Request failed:", err)
	}

	assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Response should be 201")

	assert.Equal(t, nil, tests.CleanUp(), "Cleanup test")
}

func Test_Login(t *testing.T) {
	app := tests.NewApp("../..")

	// Connect to a test db
	database.MustConnectTestDB()

	app.Post("/", requests.Use[requests.Login](), Auth.Login)

	data, _ := json.Marshal(fiber.Map{
		"email":    "jane@example.com",
		"password": "secret1234",
	})

	req := httptest.NewRequest("POST", "/", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	if err != nil {
		t.Fatal("HTTP Request failed:", err)
	}

	assert.Equal(t, fiber.StatusForbidden, res.StatusCode, "Response should be 200")

	assert.Equal(t, nil, tests.CleanUp(), "Cleanup test")
}
