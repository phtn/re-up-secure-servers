package api

import (
	"context"
	"encoding/json"
	"fast/config"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/service"
	"fast/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

const (
	Base = "/v1"
	Auth = Base + "/auth"
)

const (
	AuthRootPath      = Auth
	GetUserPath       = Auth + "/getUser"
	CreateTokenPath   = Auth + "/createToken"
	VerifyIdTokenPath = Auth + "/verifyIdToken"
	VerifyAuthKeyPath = Auth + "/verifyAuthKey"
	DevSetPath        = Auth + "/devSet"
	DevGetPath        = Auth + "/devGet"
)

var (
	conf   = config.LoadConfig()
	secret = []byte(conf.JwtSecret)
)

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.Warn("authm", "x-api-key", "missing from header")
			http.Error(w, "X-API-Key header is required", http.StatusUnauthorized)
			return
		}

		acc, err := models.GetAccountAPIKey(apiKey)
		if err != nil || acc == nil {
			utils.ErrLog("authm", "account", err)
			http.Error(w, "Invalid api key", http.StatusUnauthorized)
			return
		}
		utils.Ok("authm", "x-api-key", "matched")
		next.ServeHTTP(w, r)
	}
}

func TokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// You can add more claims verification here
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			// Add the claims to the request context for use in handlers
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
	}
}

func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if parts[1] != conf.ApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-API-Key, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	}
}

func Rdbc(w http.ResponseWriter, r *http.Request) {
	turso := models.Ping()
	redis := rdb.Ping()
	utils.JsonResponse(w, map[string]interface{}{
		"turso": turso, "redis": redis,
	})
}

func VerifyIdToken(w http.ResponseWriter, r *http.Request) {

	utils.PostMethodOnly(w, r)

	token := jwt.New(jwt.SigningMethodHS256)
	singed, err := token.SignedString(secret)
	utils.HttpError(w, "Could not generate token", err)
	utils.ErrLog("req", "verifyIdToken", err)

	utils.Ok("auth", "signed", singed)

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("req", "verifyIdToken", err)

	var v models.VerifyToken
	err = json.Unmarshal(body, &v)
	utils.ErrLog("json", "verifyIdToken", err)

	utils.JsonResponse(w, service.VerifyIdToken(r.Context(), conf.Fire, v))
}

func DevSet(w http.ResponseWriter, r *http.Request) {
	dev := "dev"

	body, err := io.ReadAll(r.Body)
	utils.ErrLog(dev, "req-body", err)

	var v auth.Token
	err = json.Unmarshal(body, &v)
	utils.ErrLog(dev, "json-unmarshal", err)

	utils.Ok(dev, "rbd-set", time.Nanosecond)
	utils.JsonResponse(w, rdb.DevSet(devKey, v))
}

func VerifyAuthKey(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("req", "verifyAuthKey", err)

	var v models.VerifyWithAuthKey
	err = json.Unmarshal(body, &v)
	utils.ErrLog("json", "verifyAuthKey", err)

	utils.JsonResponse(w, service.VerifyAuthKey(r.Context(), conf.Fire, v))
}

func CreateToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uid models.Uid
	err = json.Unmarshal(body, &uid)
	if err != nil {
		log.Fatal(err)
	}

	var response = service.CreateToken(uid, r.Context(), conf.Fire)
	utils.JsonResponse(w, response)
	log.Println(response)

}

func GetUser(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uid models.Uid
	err = json.Unmarshal(body, &uid)

	if err != nil {
		log.Fatal(err)
	}

	var response = service.GetUser(r.Context(), conf.Fire, uid)
	utils.JsonResponse(w, response)
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "retrieveToken", err)

	var k models.AuthKey
	err = json.Unmarshal(body, &k)
	utils.ErrLog("json", "retrieveToken", err)

	result, err := rdb.RetrieveToken(k.FastAuthKey)
	utils.ErrLog("redis", "retrieveToken", err)

	utils.JsonResponse(w, result)
}

func StoreVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "storeVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.ErrLog("json", "storeVal", err)

	utils.JsonResponse(w, rdb.StoreVal(kv.Key, kv.Value))
}

func RetrieveVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "retrieveVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.ErrLog("json", "retrieveVal", err)

	utils.JsonResponse(w, rdb.RetrieveVal(kv.Key))
}

func XXX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kamusta ka?")
}

var pass = os.Getenv("RDB_PASS")[:6]
var devKey = "__DEV__" + pass

func DevGet(w http.ResponseWriter, r *http.Request) {
	dev := "dev"

	_, err := io.ReadAll(r.Body)
	utils.ErrLog(dev, "req-body", err)

	utils.Ok(dev, "rbd-get", time.Nanosecond)
	utils.JsonResponse(w, rdb.DevGet(devKey))
}
