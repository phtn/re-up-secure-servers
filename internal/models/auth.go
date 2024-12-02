package models

import (
	"firebase.google.com/go/v4/auth"
)

type Uid struct {
	UID string `json:"uid"`
}

type UserAccess struct {
	AccessToken string `json:"access_token"`
	UID         string `json:"uid"`
	Email       string `json:"email"`
}

type TokenResponse struct {
	Token    *auth.Token `json:"token"`
	Error    string      `json:"error,omitempty"`
	Verified bool        `json:"verified"`
}

type UserActivation struct {
	IDToken string `json:"id_token"`
	UID     string `json:"uid"`
	Email   string `json:"email,omitempty"`
	HCode   string `json:"hcode"`
}

type ActivationResponse struct {
	GroupCode string `json:"group_code"`
	Valid     bool   `json:"valid,omitempty"`
}

type VerifyToken struct {
	IDToken   string `json:"id_token"`
	UID       string `json:"uid"`
	Email     string `json:"email"`
	GroupCode string `json:"group_code,omitempty"`
	Refresh   string `json:"refresh,omitempty"`
}
type UserCredentials struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}

type GetUserInfo struct {
	IDToken string `json:"id_token"`
	UID     string `json:"uid"`
}

type VerifyWithAuthKey struct {
	AuthKey string `json:"auth_key"`
	IDToken string `json:"id_token"`
	UID     string `json:"uid"`
	Email   string `json:"email"`
}

type AuthKey struct {
	FastAuthKey string `json:"fastinsure--auth-key"`
	DevAuthKey  string `json:"fastdev--auth-key"`
}

type VResult struct {
	IDToken   string `json:"id_token,omitempty"`
	UID       string `json:"uid,omitempty"`
	Expiry    int16  `json:"expiry,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
	GroupCode string `json:"group_code,omitempty"`
}

type V interface{}

type KV struct {
	Key   string `json:"key"`
	Value V      `json:"value"`
}

type UserCredential struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserVerified struct {
	UID      string                 `json:"uid,omitempty"`
	Claims   map[string]interface{} `json:"claims,omitempty"`
	Verified bool                   `json:"verified,omitempty"`
}

type Identity struct {
	Email []string `json:"email"`
}

type Verified struct {
	UserRecord *auth.UserRecord `json:"token,omitempty"`
	Verified   bool             `json:"verified"`
}

type Provider struct {
	SignInProvider string   `json:"sign_in_provider,omitempty"`
	Tenant         string   `json:"tenant,omitempty"`
	Identities     Identity `json:"identities,omitempty"`
}

type VerificationToken struct {
	AuthTime int64    `json:"auth_time"`
	ISS      string   `json:"iss"`
	AUD      string   `json:"aud"`
	EXP      int64    `json:"exp"`
	IAT      int64    `json:"iat"`
	SUB      string   `json:"sub,omitempty"`
	UID      string   `json:"uid,omitempty"`
	Firebase Provider `json:"firebase,omitempty"`
}
