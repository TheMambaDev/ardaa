package handlers

import (
	"ardaa/web/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	app := fiber.New()

	r := routes.NewRouter(nil, app)
	r.Setup() // change this to be a specific handler instead of loading all the routes

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/user", nil)

	res, err := app.Test(req)
	if err != nil {
		t.Errorf("an error happened %v", err)
	}

	assert.Equal(t, 200, res.StatusCode)
}
