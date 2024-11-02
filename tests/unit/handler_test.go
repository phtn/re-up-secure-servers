package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestServerLivez(t *testing.T) {

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "get http status 200",
			route:        "/",
			expectedCode: 200,
		}, {
			description:  "get http status 404",
			route:        "/not-found",
			expectedCode: 404,
		},
	}

	fibr := fiber.New()

	fibr.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("pink")
	})

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)

		resp, _ := fibr.Test(req, 1)

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
