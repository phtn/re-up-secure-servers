package api

import (
	"context"
	"fast/internal/models"
	"fast/internal/psql"
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
	L.Info("auth-mid", "check-api-key", api_key)
	if api_key == "" {
		return utils.FiberResponse(c, utils.Unauthorized, nil, utils.JsonData{Data: "api-key-missing"})
	}
	a := models.Account{}
	if err := c.BodyParser(&a); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "unable-to-parse-body"})
	}
	L.Info("auth-mid", "check-account", a)

	active, err := psql.CheckAPIKey(api_key)
	L.FailR("mdware", "api-key", err)

	result, err := service.GetUserRecordByUID(context.Background(), a.UID)
	L.Fail("get-user", "by-id", err)

	if !active || !result.Verified {
		return utils.FiberResponse(c, utils.Unauthorized, err, utils.JsonData{Data: "Unauthorized"})
	}

	L.Good("user", "active:", active, "verified:", result.Verified, err)
	return c.Next()
}

func ClaimsMiddleware(c *fiber.Ctx) error {

	var out *models.VerifyToken
	if err := c.BodyParser(&out); err != nil {
		L.Fail("mdware", "body-parser", err)
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, out)

	withClaims := t.UserRecord.CustomClaims["admin"] != nil || t.UserRecord.CustomClaims["manager"] != nil

	if t.Verified && withClaims {
		return c.Next()
	} else {
		return utils.FiberResponse(c, utils.Unauthorized, nil, utils.JsonData{Data: "Unauthorized"})
	}
}

func AdminClaimsMiddleware(c *fiber.Ctx) error {

	out := models.VerifyToken{}
	if err := c.BodyParser(out); err != nil {
		L.Fail("mdware", "body-parser", err)
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, &out)

	withClaims := t.UserRecord.CustomClaims["admin"] != nil

	if t.Verified && withClaims {
		return c.Next()
	} else {
		return utils.FiberResponse(c, utils.Unauthorized, nil, utils.JsonData{Data: "Unauthorized"})
	}
}
