package api

import (
	"context"
	"fast/internal/models"
	"fast/internal/psql"
	"fast/internal/rdb"
	"fast/internal/service"
	"fast/pkg/utils"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

type Fiddleware func(c *fiber.Ctx) fiber.Handler

var (
	L = utils.NewConsole()
	m = utils.Ice(" 𝚳 ", 0)
)

func RootRoute(c *fiber.Ctx) error {
	return OkResponse(c, "data", "OK")
}

func AuthMiddleware(c *fiber.Ctx) error {

	api_key := c.Get("X-API-Key")
	refresh_token := c.Get("X-Refresh-Token")
	auth_header := c.Get("Authorization")

	is_active, err := psql.CheckAPIKey(api_key)
	L.Fail(m, "api-key account-active", err)

	if !is_active {
		return ErrResponse(c, ErrLocked, err)
	}

	var u *rdb.Tokens
	if err := c.BodyParser(&u); err != nil {
		L.Fail(m, "body-parser", err)
		return ErrResponse(c, ErrBadRequest, err)
	}

	L.Debug(m, "UID check", u.UID)

	if u.UID == "" {
		L.Fail(m, "uid-is-empty", err)
		return ErrResponse(c, ErrUnauthorized, err)
	}

	key := "usr::" + u.UID + "::token"
	if _, ok := rdb.Int_Token_Get(key); ok {
		L.Warn(m, "key-not-found ", key)
		new_store := rdb.Int_Token_Set(key, u)
		L.Good(m, "new-key-stored", new_store)
	}

	if refresh_token == "" {
		refresh_token = u.Refresh
	}
	L.Info(m, h, "refresh-token check", refresh_token != "")
	L.Info(m, h, "refresh-token check", u.Refresh != "")

	idToken := ""
	if auth_header == "" {
		L.Warn(m, "auth-header-is-empty", err)
		idToken = u.IDToken
	} else {
		bearerToken := strings.Split(auth_header, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			return ErrResponse(c, ErrUnauthorized, err)
		}
	}

	t, err := fire.VerifyIDToken(c.Context(), u.IDToken)
	if err != nil {
		if auth.IsIDTokenExpired(err) {
			n, err := fire.VerifyIDTokenAndCheckRevoked(c.Context(), u.Refresh)
			if err != nil {
				L.Fail(m, "expired-get-token-with-refresh", err)
				return ErrResponse(c, ErrUnauthorized, err)
			}

			c.Locals("id_token", n)
			c.Locals("refresh_token", u.Refresh)
			c.Locals("auth_token", n)

			return c.Next()

		}
		return ErrResponse(c, ErrUnauthorized, err)
	}

	c.Locals("id_token", idToken)
	c.Locals("refresh_token", refresh_token)
	c.Locals("auth_token", t)

	L.Info(m, "validated", t.UID)

	return c.Next()

}

func AdminClaimsMiddleware(c *fiber.Ctx) error {

	out := models.VerifyToken{}
	if err := c.BodyParser(out); err != nil {
		L.Fail(m, "body-parser", err)
		return ErrResponse(c, ErrBadRequest, err)
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, &out)

	withClaims := t.UserRecord.CustomClaims["admin"] != nil

	if t.Verified && withClaims {
		return c.Next()
	} else {
		return ErrResponse(c, ErrUnauthorized, nil)
	}
}

func ClaimsMiddleware(c *fiber.Ctx) error {

	var out *models.VerifyToken
	if err := c.BodyParser(&out); err != nil {
		L.Fail(m, "body-parser", err)
		return ErrResponse(c, ErrBadRequest, err)
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, out)

	withClaims := t.UserRecord.CustomClaims["admin"] != nil || t.UserRecord.CustomClaims["manager"] != nil

	if t.Verified && withClaims {
		return c.Next()
	} else {
		return ErrResponse(c, ErrUnauthorized, nil)
	}
}
