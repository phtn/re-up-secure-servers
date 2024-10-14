package models

import "github.com/golang-jwt/jwt/v5"

type Account struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UID   string `json:"uid"`
	Email string `json:"email"`
}
