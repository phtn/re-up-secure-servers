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
	AuthPath          = Auth
	GetUserPath       = "/get-user"
	CreateTokenPath   = "/create-token"
	VerifyIdTokenPath = "/verify-id-token"
	VerifyAuthKeyPath = "/verify-auth-key"
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
	out := new(models.VerifyToken)
	if err := c.BodyParser(out); err != nil {
		return utils.FiberResponse(c, utils.BadRequest, err, utils.JsonData{Data: out})
	}
	d := service.VerifyIdToken(c.Context(), fire, out)

	data := utils.JsonData{Data: d}
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
	t := service.GetUserRecord(ctx, fire, out)

	data := utils.JsonData{Data: t.UserRecord.CustomClaims}
	L.Info("get claims", "data", data)

	return utils.FiberResponse(c, utils.OK, nil, data)
}
