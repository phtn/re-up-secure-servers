package models

import "firebase.google.com/go/v4/auth"

type Uid struct {
	UID string `json:"uid"`
}

type UserAccess struct {
	AccessToken string `json:"access_token"`
	UID         string `json:"uid"`
	Email       string `json:"email"`
}

type UserRefresh struct {
	RefreshToken string `json:"refresh_token"`
	UID          string `json:"uid"`
	Email        string `json:"email"`
}

type VerifyToken struct {
	IDToken   string `json:"id_token"`
	UID       string `json:"uid"`
	Email     string `json:"email"`
	GroupCode string `json:"group_code,omitempty"`
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
	Key      string `json:"key"`
	Verified bool   `json:"verified"`
	Exp      int16  `json:"exp"`
	IsActive bool   `josn:"is_active,omitempty"`
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
