package models

type Uid struct {
	UID string `json:"uid"`
}

type VerifyToken struct {
	IDToken string `json:"idToken"`
	UID     string `json:"uid"`
	Email   string `json:"email"`
}

type VerifyWithAuthKey struct {
	AuthKey string `json:"authKey"`
	IDToken string `json:"idToken"`
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
