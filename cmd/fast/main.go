package main

import (
	"fast/api"
	"fast/config"
	"fast/internal/psql"
	"fast/internal/rdb"
	"fast/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
)

func main() {

	addr := config.LoadConfig().Addr

	server := fiber.New()
	server.Use(idempotency.New(idempotency.Config{
		Lifetime: 30 * time.Minute,
	}))
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// server.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-CSRF-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.Guid,
	// }))

	withAuth := []fiber.Handler{api.AuthMiddleware}
	withClaims := []fiber.Handler{api.AuthMiddleware, api.ClaimsMiddleware}
	withAdmin := []fiber.Handler{api.AuthMiddleware, api.ClaimsMiddleware, api.AdminClaimsMiddleware}

	// SERVER
	F := server.Group("/")
	F.Get(api.Livez, api.ServerLivez)
	F.Get(api.Readyz, api.ServerReadyz)
	F.Post(api.VerifyAgentCodePath, api.VerifyAgentCode)

	// AUTHENTICATED
	authGroup := server.Group(api.AuthPath, withAuth...)
	authGroup.Post(api.VerifyIdTokenPath, api.VerifyUser)
	authGroup.Post(api.VerifyUserPath, api.VerifyUser)
	authGroup.Post(api.VerifyOnSigninPath, api.GetUser)
	authGroup.Post(api.GetUserPath, api.GetUserInfo)
	authGroup.Post(api.GetClaimsPath, api.GetClaims)
	authGroup.Post(api.ActivateUserPath, api.ActivateUser)

	// WITH CLAIMS
	claimsGroup := server.Group(api.ClaimsPath, withClaims...)
	claimsGroup.Post(api.CustomClaimsPath, withClaims...)
	claimsGroup.Post(api.AgentCodePath, api.CreateAgentCode)

	// WITH ADMIN CLAIMS
	adminGroup := server.Group(api.AdminPath, withAdmin...)
	adminGroup.Post(api.AdminPath, withAdmin...)
	adminGroup.Post(api.AccountTokenPath, api.CreateAccountToken)

	// DEV
	debugGroup := server.Group(api.DevPath)
	debugGroup.Post(api.DebugRedisPath, api.DebugRedisStore)

	utils.MkOne()

	// DB HEALTH //
	psql.PsqlHealth()
	rdb.RedisHealth()
	// END DB HEALTH //

	// TEST //
	//
	// customClaims := service.CustomClaims{
	// 	"agent":   "true",
	// 	"manager": "true",
	// }
	//

	// END TEST //

	// SERVER START
	server.Listen(addr)
}
