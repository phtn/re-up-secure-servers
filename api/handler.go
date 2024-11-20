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
)

const (
	AuthPath            = Auth
	GetUserPath         = "/get-user"
	CreateTokenPath     = "/create-token"
	VerifyIdTokenPath   = "/verify-id-token"
	VerifyAuthKeyPath   = "/verify-auth-key"
	VerifyAgentCodePath = "/verify-agent-code"
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
)

var (
	fire   = config.LoadConfig().Fire
	secret = []byte(config.LoadConfig().JwtSecret)
)

func ServerLivez(c *fiber.Ctx) error {
	data := utils.JsonData{Data: "OK"}
	L.Info("server  ", "livez", data)
	return utils.FiberResponse(c, utils.OK, nil, data)
}

func ServerReadyz(c *fiber.Ctx) error {
	data := utils.JsonData{Data: "OK"}
	L.Info("server  ", "readyz", data)
	return utils.FiberResponse(c, utils.OK, nil, data)
}

func DatabaseHealth(c *fiber.Ctx) error {
	psql := psql.PsqlHealth()
	rdbs := rdb.RedisHealth()
	data := Health{
		Psql: psql, Rdbs: rdbs,
	}
	return utils.FiberResponse(c, utils.OK, nil, utils.JsonData{Data: data})
}

func VerifyIdToken(c *fiber.Ctx) error {
	var out models.VerifyToken
	if err := c.BodyParser(&out); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: out})
	}

	result := service.VerifyIdToken(c.Context(), out)

	data := utils.JsonData{Data: result}
	return utils.FiberCookie(c, result.Cookie, utils.OK, nil, data)
}

func GetUserInfo(c *fiber.Ctx) error {
	var out models.GetUserInfo
	var data utils.JsonData
	if err := c.BodyParser(&out); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: out})
	}
	refresh := c.Get("x-refresh-token")
	user_refresh := models.UserRefresh{
		IDToken: out.IDToken,
		UID:     out.UID,
		Refresh: refresh,
	}
	L.Info("INFO", "handler: GetUserInfo", user_refresh.Refresh)

	response, err := service.TokenVerification(c.Context(), user_refresh)
	L.Warn("verification", "service: TokenVerification", err)
	if err != nil {
		return utils.FiberResponse(c, utils.Unauthorized, err, data)
	}

	data = utils.JsonData{Data: &response.Verified}
	if !response.Verified {
		return utils.FiberResponse(c, utils.Unauthorized, nil, data)
	}

	user, err := service.GetUserInfo(c.Context(), out.UID)
	L.FailR("get-user", "firebase", err)

	data = utils.JsonData{Data: user}
	return utils.FiberResponse(c, utils.OK, nil, data)
}

func CreateAgentCode(c *fiber.Ctx) error {
	var v models.VerifyToken
	if err := c.BodyParser(&v); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	result := service.NewAgentCode(v)
	data := utils.JsonData{Data: result}
	return utils.FiberResponse(c, utils.OK, nil, data)
}

func VerifyAgentCode(c *fiber.Ctx) error {
	var p *models.HCodeParams
	if err := c.BodyParser(&p); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	result := service.VerifyAgentCode(p)
	data := utils.JsonData{Data: result}
	return utils.FiberResponse(c, utils.OK, nil, data)
}

func CreateAccountToken(c *fiber.Ctx) error {
	var v models.VerifyToken
	if err := c.BodyParser(&v); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	var u = shield.NewAccountToken{UID: v.UID, Email: v.Email}
	data, err := shield.NewAccount(&u)
	if err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	return utils.FiberResponse(c, utils.OK, nil, utils.JsonData{Message: "  OK - POST", Data: data})
}

func GetClaims(c *fiber.Ctx) error {
	out := new(models.VerifyToken)
	if err := c.BodyParser(out); err != nil {
		L.Fail("mdware", "body-parser", err)
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: "Bad Request", Error: err, Message: "body-params-invalid"})
	}
	ctx := context.Background()
	t := service.GetUserRecord(ctx, out)

	data := utils.JsonData{Data: t.UserRecord}
	L.Info("get claims", "data", data)

	return utils.FiberResponse(c, utils.OK, nil, data)
}
