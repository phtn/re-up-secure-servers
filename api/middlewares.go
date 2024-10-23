package api

import (
	"context"
	"fast/internal/models"
	"fast/internal/service"
	"fast/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type Fiddleware func(c *fiber.Ctx) fiber.Handler

var (
	L = utils.NewConsole()
)

func RootRoute(c *fiber.Ctx) error {
	return utils.FiberResponse(c, utils.OK, nil, utils.JsonData{Data: "OK"})
}

func AuthMiddleware(c *fiber.Ctx) error {
	api_key := c.Get("X-API-Key")
	if api_key == "" {
		return utils.FiberResponse(c, utils.Unauthorized, nil, utils.JsonData{Data: "api-key-missing"})
	}
	a := new(models.Account)
	if err := c.BodyParser(&a); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "unable-to-parse-body"})
	}
	account, err := models.GetAccountWithAPIKey(api_key)
	L.Fail("mdware", "api-key", err)

	if !account.Active {
		return utils.FiberResponse(c, utils.Unauthorized, err, utils.JsonData{Data: "Unauthorized"})
	}

	return c.Next()
}

func ClaimsMiddleware(c *fiber.Ctx) error {

	out := new(models.VerifyToken)
	if err := c.BodyParser(out); err != nil {
		L.Fail("mdware", "body-parser", err)
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, fire, out)

	withClaims := t.UserRecord.CustomClaims["admin"] != nil || t.UserRecord.CustomClaims["manager"] != nil

	if t.Verified && withClaims {
		return c.Next()
	} else {
		return utils.FiberResponse(c, utils.Unauthorized, nil, utils.JsonData{Data: "Unauthorized"})
	}
}
