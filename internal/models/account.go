package models

import (
	"context"
	"encoding/json"
	"fast/config"
	"fast/internal/shield"
	"fast/pkg/utils"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	L = utils.NewConsole()
)

type Role string

const (
	Dev   Role = "developer"
	Adm   Role = "admin"
	Mgr   Role = "manager"
	Acm   Role = "account_manager"
	Agent Role = "agent"
	Suprv Role = "supervisor"
	Undrw Role = "underwriter"
	Dealr Role = "dealer"
)

type Account struct {
	UID    string `json:"uid,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	APIKey string `json:"api_key,omitempty"`
	Active bool   `json:"is_active,omitempty"`
	Role   Role   `json:"role,omitempty"`
	CrTime string `json:"creation_time,omitempty"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UID   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role  Role   `json:"role,omitempty"`
}

var (
	jwts = []byte(config.LoadConfig().JwtSecret)
	fire = config.LoadConfig().Fire
	rdb  = config.LoadConfig().Rdbs
)

func NewAccountCustomClaims(acct Account) (string, error) {

	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "re-up.ph secure servers",
			Subject:   acct.UID,
			Audience:  []string{"re-up secure servers clients"},
		},
		UID:  acct.UID,
		Role: acct.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwts)
	L.Fail("claim", "unable to signed token", err)

	return verifyToken(tokenString)
}

func verifyToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwts, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		storeClaims(shield.NewClaimsKey(claims.Email), claims)
	}

	return "", fmt.Errorf("invalid token")
}

func storeClaims(key string, v *CustomClaims) {

	value, err := json.Marshal(&v)
	L.Fail("json", "marshal-store-claims", err)

	ctx := context.Background()
	err = rdb.Set(ctx, key, value, 24*7*time.Hour).Err()
	L.Fail("redis", "set store-claims", err)

	L.Good("redis", "set complete", key, err)
}

func GenerateRefreshToken(user Account, key []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.UID,
		Issuer:    "re-up.ph secure servers",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		Audience:  []string{"Re-up Secure Servers Clients"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func GenerateAccessToken(user Account) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.UID,
		Issuer:    "re-up.ph secure servers",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Audience:  []string{"Re-up Secure Servers Clients"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwts)
}
