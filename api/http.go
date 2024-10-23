package api

import (
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/golang-jwt/jwt/v4"
)

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func _DbCheck(w http.ResponseWriter, r *http.Request) {
	psql := models.PsqlHealth()
	redis := rdb.RedisHealth()
	utils.JsonResponse(w, map[string]interface{}{
		"psql": psql, "redis": redis,
	})
}

func _VerifyIdToken(w http.ResponseWriter, r *http.Request) {

	utils.PostMethodOnly(w, r)

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("req", "verifyIdToken", err)

	var v models.VerifyToken
	err = json.Unmarshal(body, &v)
	utils.ErrLog("json", "verifyIdToken", err)

	// utils.JsonResponse(w, service.VerifyIdToken(r.Context(), fire, v))
}

func _VerifyAuthKey(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body-params", err, http.StatusBadRequest)

	var v models.VerifyWithAuthKey
	err = json.Unmarshal(body, &v)
	utils.HttpError(w, "json-body", err, http.StatusBadRequest)

	utils.JsonResponse(w, service.VerifyAuthKey(r.Context(), fire, v))
}

func _CreateCustomClaims(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var u *service.UserCredentials
	err = json.Unmarshal(body, &u)
	utils.HttpError(w, "Unable to read json", err, http.StatusBadRequest)

	data, err := service.NewCustomClaims(u)
	utils.HttpError(w, "Unable to read json", err, http.StatusServiceUnavailable)

	utils.JsonResponse(w, data)
}

func _CheckAdminAuthority(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var v *service.UserCredentials
	err = json.Unmarshal(body, &v)
	if err != nil {
		utils.Info("handl", "admin-authority", v)
		utils.ErrLog("handl", "admin-authority", err)
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	verified := service.VerifyAdmin(r.Context(), fire, v)
	utils.HttpError(w, "unable to verify credentials", err, http.StatusServiceUnavailable)

	if !verified {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	data := map[string]interface{}{
		"verified":  verified,
		"timestamp": time.Now().Unix(),
	}

	utils.JsonResponse(w, data)
}

func _CreateAdminClaims(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var u *service.UserCredentials
	err = json.Unmarshal(body, &u)
	utils.HttpError(w, "Unable to read json", err, http.StatusBadRequest)

	data, err := service.NewAdminCustomClaims(u)
	utils.HttpError(w, "Unable to read json", err, http.StatusServiceUnavailable)

	utils.JsonResponse(w, data)
}

func _CreateAgentCode(w http.ResponseWriter, r *http.Request) {

	utils.PostMethodOnly(w, r)
	url := r.URL.Query().Get("url")

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("acode", "create", err)

	var v models.VerifyToken
	err = json.Unmarshal(body, &v)
	utils.ErrLog("acode", "jsond", err)
	utils.Info("handl", "url", url)

	data := service.NewAgentCode(v)
	utils.JsonResponse(w, data)
}

func _CreateOneTimeClaim(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var u *service.UserCredentials
	err = json.Unmarshal(body, &u)
	utils.HttpError(w, "Unable to read json", err, http.StatusBadRequest)

	data, err := service.NewOneTimeClaim(u)
	utils.HttpError(w, "Unable to read json", err, http.StatusServiceUnavailable)

	utils.JsonResponse(w, data)
}

func _CreateToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uid models.Uid
	err = json.Unmarshal(body, &uid)
	if err != nil {
		log.Fatal(err)
	}

	var response = service.CreateToken(uid, r.Context(), fire)
	utils.JsonResponse(w, response)
	log.Println(response)

}

func _GetUser(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uid models.Uid
	err = json.Unmarshal(body, &uid)

	if err != nil {
		log.Fatal(err)
	}

	var response = service.GetUser(r.Context(), fire, uid)
	utils.JsonResponse(w, response)
}

func _RetrieveToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "retrieveToken", err)

	var k models.AuthKey
	err = json.Unmarshal(body, &k)
	utils.ErrLog("json", "retrieveToken", err)

	result, err := rdb.RetrieveToken(k.FastAuthKey)
	utils.ErrLog("redis", "retrieveToken", err)

	utils.JsonResponse(w, result)
}

func _StoreVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "storeVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.ErrLog("json", "storeVal", err)

	utils.JsonResponse(w, rdb.StoreVal(kv.Key, 24, kv.Value))
}

func _RetrieveVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "retrieveVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.ErrLog("json", "retrieveVal", err)

	utils.JsonResponse(w, rdb.RetrieveVal(kv.Key))
}

func _XXX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kamusta ka?")
}

var pass = os.Getenv("RDB_PASS")[:6]
var devKey = "__DEV__" + pass

func _DevGet(w http.ResponseWriter, r *http.Request) {
	dev := "dev"

	_, err := io.ReadAll(r.Body)
	utils.ErrLog(dev, "req-body", err)

	utils.Ok(dev, "rbd-get", time.Nanosecond)
	utils.JsonResponse(w, rdb.DevGet(devKey))
}

func _DevSet(w http.ResponseWriter, r *http.Request) {
	dev := "dev"

	body, err := io.ReadAll(r.Body)
	utils.ErrLog(dev, "req-body", err)

	var v auth.Token
	err = json.Unmarshal(body, &v)
	utils.ErrLog(dev, "json-unmarshal", err)

	utils.Ok(dev, "rbd-set", time.Nanosecond)
	utils.JsonResponse(w, rdb.DevSet(devKey, v))
}

// token := jwt.New(jwt.SigningMethodHS256)
// singed, err := token.SignedString(secret)
// utils.HttpError(w, "Could not generate token", err, http.StatusInternalServerError)
// utils.ErrLog("req", "verifyIdToken", err)

// utils.Ok("auth", "signed", singed)

func _AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject := r.UserAgent() + " " + r.RemoteAddr + " " + r.Referer()
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.Warn("authm", "x-api-key missing", subject)
			http.Error(w, "X-API-Key header is required", http.StatusUnauthorized)
			return
		}

		a, err := models.GetAccountWithAPIKey(apiKey)
		utils.HttpError(w, "api-key", err, http.StatusUnauthorized)

		if !a.Active {
			utils.HttpError(w, "account-not-active", err, http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	}
}

func _ClaimsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := io.ReadAll(r.Body)
		utils.ErrLog("req_b", "verifyIdToken", err)

		var v models.VerifyToken
		err = json.Unmarshal(data, &v)
		utils.ErrLog("jsond", "body-params", err)

		ctx := context.Background()
		t := service.GetUserRecord(ctx, fire, &v)
		utils.ErrLog("authv", "unable to verify id token", err)

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		withClaims := t.UserRecord.CustomClaims["admin"] != nil || t.UserRecord.CustomClaims["manager"] != nil

		utils.Ok("mdlwr", "verified", t.Verified)
		utils.Info("mdlwr", "claims", withClaims)

		if !t.Verified && !withClaims {
			utils.HttpErr(w, "Unauthorized", err, http.StatusUnauthorized)
		}

		if t.Verified && withClaims {
			next.ServeHTTP(w, r)
		}
	}
}

func AdminClaimsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.Warn("authm", "x-api-key missing", "api-key-check")
			http.Error(w, "X-API-Key header is required", http.StatusUnauthorized)
			return
		}

		data, err := io.ReadAll(r.Body)
		utils.ErrLog("req_b", "verifyIdToken", err)

		var v *service.UserCredentials
		err = json.Unmarshal(data, &v)
		utils.ErrLog("jsond", "unable to unmarshal body contents", err)
		// utils.Info("jsond", "body-contents", string(data))

		if v.Claims.Role == "" || v.Claims.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			utils.ErrLog("claim", "verification-failed", err)
		}
		utils.Info("middl", "role", v.Claims.Role)

		r.Body = io.NopCloser(bytes.NewBuffer(data))

		next.ServeHTTP(w, r)
	}
}

func apiKeyCheck(apiKey string, w http.ResponseWriter, r *http.Request, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.Warn("authm", "x-api-key missing", "api-key-check")
			http.Error(w, "X-API-Key header is required", http.StatusUnauthorized)
			return
		}
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

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-API-Key, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
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

		if parts[1] != "api_key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
