package api

import (
	"context"
	"fast/config"
	"fast/internal/models"
	"fast/internal/psql"
	"fast/internal/rdb"
	"fast/internal/service"
	"fast/internal/shield"
	"fast/pkg/utils"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

type Health struct {
	Psql interface{} `json:"psql,omitempty"`
	Rdbs interface{} `json:"rdbs,omitempty"`
}

const (
	Root   = "/v1"
	Livez  = "livez"
	Readyz = "readyz"
	Auth   = Root + "/auth"
	Admin  = Root + "/admin"
	Claims = Root + "/claims"
	Dev    = Root + "/dev"
)

const (
	AuthPath            = Auth
	GetUserPath         = "/get-user"
	CreateTokenPath     = "/create-token"
	VerifyUserPath      = "/verify-user"
	VerifyIdTokenPath   = "/verify-id-token"
	VerifyOnSigninPath  = "/verify-on-signin"
	VerifyAuthKeyPath   = "/verify-auth-key"
	VerifyAgentCodePath = "/verify-agent-code"
	ActivateUserPath    = "/activate-user"
	// CLAIMS
	ClaimsPath       = Claims
	CustomClaimsPath = "/create-custom-claims"
	GetClaimsPath    = "/get-claims"
	AgentCodePath    = "/create-code"
	CreateGroupPath  = "/create-group"
	// DEV
	DevSetPath = "/dev-set"
	DevGetPath = "/dev-get"
	// ADMIN
	AdminPath        = Admin
	AdminClaimsPath  = "/admin-claims"
	AccountTokenPath = "/create-account-token"
	DevPath          = Dev
	DebugRedisPath   = "/rdb-debug"
)

const (
	OK           = fiber.StatusOK
	Accepted     = fiber.StatusAccepted
	Unauthorized = fiber.StatusUnauthorized
	BadRequest   = fiber.StatusBadRequest
	NotFound     = fiber.StatusNotFound
	Processing   = fiber.StatusProcessing
)

var (
	fire   = config.LoadConfig().Fire.AuthClient
	secret = []byte(config.LoadConfig().JwtSecret)
	h      = utils.Sky("𝚑𝚝𝚝𝚙", 0)
)

func ServerLivez(c *fiber.Ctx) error {
	L.Info(h, "livez", "OK")
	return OkResponse(c, "All systems good.", "⌨")
}

func ServerReadyz(c *fiber.Ctx) error {
	L.Info(h, "readyz", "OK")
	return OkResponse(c, "All systems ready", "⚡")
}

func DatabaseHealth(c *fiber.Ctx) error {
	psql := psql.PsqlHealth()
	rdbs := rdb.RedisHealth()
	data := Health{
		Psql: psql, Rdbs: rdbs,
	}
	return OkResponse(c, Health{data, ""}, nil)
}

func VerifyUser(c *fiber.Ctx) error {

	id_token := c.Locals("id_token").(string)
	refresh_token := c.Locals("refresh_token").(string)
	auth_token := c.Locals("auth_token").(*auth.Token)

	var (
		expiresIn = 24 * 5 * time.Hour
	)
	u := rdb.Tokens{
		IDToken: id_token,
		Refresh: refresh_token,
		UID:     auth_token.UID,
		Expiry:  time.Now().Add(expiresIn),
	}

	errors := ValidateFields(u)

	if errors != nil {
		v_errors := AppError{
			Status:  BadRequest,
			Code:    ErrCodeBadRequest,
			Message: "Missing fields",
			Details: errors}

		return ErrResponse(c, &v_errors, nil)
	}

	// uid := "usr::" + auth_token.UID + "::token"
	// store, err := rdb.GetUserTokens(uid)
	// if err != nil {
	// return ErrResponse(c, ErrUnauthorized, err)
	// }

	// session, err := fire.VerifySessionCookieAndCheckRevoked(c.Context(), u.IDToken)
	// if err != nil {
	// 	return ErrResponse(c, ErrUnauthorized, err)
	// }
	// if time.Until(store.Expiry) < 5*time.Minute {
	// 		L.Warn(h, "expiry is less that 5 mins", time.Until(store.Expiry))
	// 		cookieStr, err = fire.SessionCookie(c.Context(), id_token, expiresIn)
	// 		ErrResponse(c, ErrUnauthorized, err)
	// 	}
	// cookie, err := fire.SessionCookie(c.Context(), id_token, expiresIn)

	// var details interface{}
	// if session.UID == u.UID {
	// 	claims := session.Claims

	// 	details = fiber.Map{
	// 		"status": OK,
	// 		"data": models.UserVerified{
	// 			UID:      session.UID,
	// 			Claims:   claims,
	// 			Verified: true,
	// 		},
	// 	}
	// }

	// L.Good(h, "key found", "OK")
	// data := DataResponse{
	// 	Status:  200,
	// 	Message: "success",
	// 	Code:    "VERIFICATION",
	// 	Err:     nil,
	// 	Data:    details,
	// }
	// return CookieHandler(c, cookie, &data)
	return OkResponse(c, "", "")
}

///////////////////////////////////////////////////

func VerifyIdToken(c *fiber.Ctx) error {
	var out models.VerifyToken
	if err := c.BodyParser(&out); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}

	result, err := service.VerifyIdToken(c.Context(), out)
	L.Fail(h, "verification", err)
	data := Res{Data: result}
	return OkResponse(c, data, nil)
}

func GetUserInfo(c *fiber.Ctx) error {
	var out models.GetUserInfo
	if err := c.BodyParser(&out); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	refresh := c.Get("x-refresh-token")
	user_refresh := rdb.UserTokens{
		IDToken: out.IDToken,
		UID:     out.UID,
		Refresh: refresh,
	}
	L.Info(h, "handler: GetUserInfo", user_refresh.Refresh)

	response, err := service.TokenVerification(c.Context(), user_refresh)
	L.Warn(h, "service: TokenVerification", err)
	if err != nil {
		return ErrResponse(c, ErrUnauthorized, err)
	}

	if !response.Verified {
		return ErrResponse(c, ErrUnauthorized, nil)
	}

	user, err := service.GetUserInfo(c.Context(), out.UID)
	L.FailR(h, "get-user-record", err)

	return OkResponse(c, user, nil)
}

func CreateAgentCode(c *fiber.Ctx) error {
	var v models.VerifyToken
	if err := c.BodyParser(&v); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	result := service.NewAgentCode(v)
	data := Res{Data: result}
	return OkResponse(c, data, nil)
}

func ActivateUser(c *fiber.Ctx) error {
	var v models.UserActivation
	if err := c.BodyParser(&v); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	result, err := service.GetUserRecordByUID(context.Background(), v.UID)
	L.Fail(h, "activate get-user", err)

	if err != nil {
		return ErrResponse(c, ErrUnauthorized, err)
	}
	if result.Verified {

		a := service.UnlockWithKey(v.HCode)
		data := Res{Data: &a}

		claim := map[string]interface{}{"agent": true}
		claims, err := service.AddCustomClaim(v.IDToken, v.UID, claim)
		L.Fail(h, "activation add-agent-claim", err)
		data = Res{Data: "Unable to add custom claims for uid"}
		if err != nil {
			return ErrResponse(c, ErrUnauthorized, err)
		}

		withClaims := claims.Claims["agent"] != nil

		if withClaims {
			L.Good(h, "activation custom-claim-added", withClaims)

			if exists, user := psql.CheckIfUserExists(v.UID); exists {
				res, err := user.Update().SetGroupCode(a.GroupCode).Save(context.Background())
				L.Fail("activate-user", "set-group-code", err)
				L.Good("activate-user", "set-group-code", res.GroupCode)

				return OkResponse(c, data, nil)
			}

			phone_number := service.MockPhone()
			new_user := psql.NewUser(v.Email, v.Email, phone_number, v.UID, a.GroupCode)

			if new_user == v.UID {
				L.Good(h, "uid", new_user)
			}
			return OkResponse(c, OK, nil)
		}
		L.Warn(h, "claims-not-set", v.UID)
	}
	return ErrResponse(c, ErrUnauthorized, err)
}

func VerifyAgentCode(c *fiber.Ctx) error {
	var p *models.HCodeParams
	if err := c.BodyParser(&p); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	result := service.VerifyAgentCode(p)
	data := Res{Data: result}
	return OkResponse(c, data, nil)
}

func CreateAccountToken(c *fiber.Ctx) error {
	var v models.VerifyToken
	if err := c.BodyParser(&v); err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	var u = shield.NewAccountToken{UID: v.UID, Email: v.Email}
	data, err := shield.NewAccount(&u)
	if err != nil {
		return ErrResponse(c, ErrBadRequest, err)
	}
	return OkResponse(c, data, nil)
}

func GetClaims(c *fiber.Ctx) error {
	out := new(models.VerifyToken)
	if err := c.BodyParser(out); err != nil {
		L.Fail(h, "get-claim body-parser", err)
		return ErrResponse(c, ErrBadRequest, err)
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, out)

	data := Res{Data: t.UserRecord}
	L.Info(h, "data", data)

	return OkResponse(c, data, nil)
}
