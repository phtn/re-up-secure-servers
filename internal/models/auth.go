package models

type Uid struct {
	UID string `json:"uid"`
}

type IdToken struct {
	Token string `json:"idToken"`
}

type UserCredential struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
