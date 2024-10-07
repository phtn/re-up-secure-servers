package models

type Uid struct {
	UID string `json:"uid"`
}

type IdToken struct {
	Token string `json:"idToken"`
}

type AuthKey struct {
	FastAuthKey string `json:"fast_auth_key"`
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
