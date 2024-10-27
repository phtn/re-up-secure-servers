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

	// AUTHENTICATED
	authGroup := server.Group(api.AuthPath, withAuth...)
	authGroup.Post(api.VerifyIdTokenPath, api.VerifyIdToken)
	authGroup.Post(api.GetClaimsPath, api.GetClaims)

	// WITH CLAIMS
	claimsGroup := server.Group(api.ClaimsPath, withClaims...)
	claimsGroup.Post(api.CustomClaimsPath, withClaims...)
	claimsGroup.Post(api.AgentCodePath, api.CreateAgentCode)

	// WITH ADMIN CLAIMS
	adminGroup := server.Group(api.AdminPath, withAdmin...)
	adminGroup.Post(api.AdminPath, withAdmin...)
	adminGroup.Post(api.AccountTokenPath, api.CreateAccountToken)

	utils.MkOne()

	// DB HEALTH //
	psql.PsqlHealth()
	rdb.RedisHealth()
	// END DB HEALTH //

	// TEST //
	// psql.CreateAccount("re-up.ph", "hq@re-up.ph", "ZjI2NTk3MjQtMzRmNi00MGFjLThhZmItYjE1OGFmYTJmNzM0", "N7yCd3kCViMA0jD3eNuv5rqKxgy1")
	// psql.GetAllAccounts()

	// END TEST //

	// SERVER START
	server.Listen(addr)
}

// mux := http.NewServeMux()
// admin_middlewares := []api.Middleware{
// 	api.AuthMiddleware,
// 	api.CorsMiddleware,
// 	api.AdminClaimsMiddleware,
// }
// withClaims := append(middlewares, api.ClaimsMiddleware)
// withAdmin := append(middlewares, api.AdminClaimsMiddleware)

// authGroup.Get(api.C)
// authGroup.Get(api.GetUserPath)
// authGroup.Post(api.CreateTokenPath)
// authGroup.Post(api.VerifyAuthKeyPath)

// mux.HandleFunc(api.AuthPath, api.Chain(api.DbCheck, middlewares...))
// mux.HandleFunc(api.GetUserPath, api.Chain(api.GetUser, middlewares...))
// mux.HandleFunc(api.CreateTokenPath, api.Chain(api.CreateToken, middlewares...))
// mux.HandleFunc(api.VerifyIdTokenPath, api.Chain(api.VerifyIdToken, middlewares...))
// mux.HandleFunc(api.VerifyAuthKeyPath, api.Chain(api.VerifyAuthKey, middlewares...))

// // WITH CLAIMS
// mux.HandleFunc(api.CustomClaimsPath, api.Chain(api.CreateCustomClaims, withClaims...))
// mux.HandleFunc(api.AgentCodePath, api.Chain(api.CreateAgentCode, withClaims...))
// // WITH ADMIN
// mux.HandleFunc(api.AdminPath, api.Chain(api.CheckAdminAuthority, admin_middlewares...))
// mux.HandleFunc(api.AdminClaimsPath, api.Chain(api.CreateAdminClaims, withAdmin...))

// // DEV-ROUTES
// mux.HandleFunc(api.DevSetPath, api.Chain(api.DevSet, middlewares...))
// mux.HandleFunc(api.DevGetPath, api.Chain(api.DevGet, middlewares...))

// server := &http.Server{
// 	Addr:    addr,
// 	Handler: mux,
// }
// err := server.ListenAndServe()
// utils.Fatal("serve", "boot", err)
// utils.OkLog("serve", "boot", "system-online", err)
