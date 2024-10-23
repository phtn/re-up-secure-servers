package api

import (
	"fast/config"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/service"
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
	GetUserPath       = "/getUser"
	CreateTokenPath   = "/createToken"
	VerifyIdTokenPath = "/verifyIdToken"
	VerifyAuthKeyPath = "/verifyAuthKey"
	// CLAIMS
	ClaimsPath       = Claims
	CustomClaimsPath = "/createCustomClaims"
	AgentCodePath    = "/create-code"
	DevSetPath       = "/devSet"
	DevGetPath       = "/devGet"
	// ADMIN
	AdminPath       = Admin
	AdminClaimsPath = "/adminClaims"
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
	psql := models.PsqlHealth()
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
	data := service.NewAgentCode(v)
	return utils.FiberResponse(c, utils.OK, nil, utils.JsonData{Message: "  OK - POST", Data: data})
}
