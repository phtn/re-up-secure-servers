package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fast/config"
	"fast/internal/models"
	"fast/internal/service"
	"fast/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	api_key        = config.LoadConfig().ApiKey
	allowed_origin = config.LoadConfig().AllowedOrigin
)

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func ClaimsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := io.ReadAll(r.Body)
		utils.ErrLog("req_b", "verifyIdToken", err)

		var v *service.UserCredentials
		err = json.Unmarshal(data, &v)
		utils.ErrLog("jsond", "body-params", err)

		if v.Claims.Role == "" || v.Claims.Role != "manager" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			utils.ErrLog("claim", "verification-failed", err)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		next.ServeHTTP(w, r)
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
		w.Header().Set("Access-Control-Allow-Origin", allowed_origin)
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

		if parts[1] != api_key {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
