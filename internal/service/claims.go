package service

import (
	"context"
	"fast/internal/psql"

	"firebase.google.com/go/v4/auth"
)

type Claims map[string]interface{}

type UserClaims struct {
	Role string `json:"role,omitempty"`
}

type UserCredentials struct {
	UID          string     `json:"uid,omitempty"`
	Email        string     `json:"email,omitempty"`
	Claims       UserClaims `json:"claims,omitempty"`
	IDToken      string     `json:"id_token,omitempty"`
	Subject      string     `json:"sub,omitempty"`
	AdminKey     string     `json:"admin_key,omitempty"`
	CustomClaims Claims     `json:"custom_claims,omitempty"`
}

var (
	ctx = context.Background()
)

func AddCustomClaim(idToken string, uid string, claim string) (*auth.Token, error) {

	customClaims := Claims{claim: true}

	err := fire.SetCustomUserClaims(ctx, uid, customClaims)
	if err != nil {
		return nil, err
	}

	L.Info("check set", uid, customClaims)

	token, err := fire.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	claims := token.Claims

	for k, value := range customClaims {
		switch v := value.(type) {
		case string:
			if withClaim, ok := claims[k]; ok {
				L.Good("claim", k, "confirmed", withClaim)
				return token, nil
			}

		default:
			L.Info("agent", "k", k, "value", v, value)
			return token, nil
		}
	}
	return nil, err
}

func NewCustomClaims(u *UserCredentials) (*auth.Token, error) {

	customClaims := u.CustomClaims
	t, err := fire.VerifyIDToken(ctx, u.IDToken)
	L.Fail("authv", "unable to verify id token", err)

	verified := t.UID == u.UID
	authorized := false

	if role, ok := t.Claims["role"]; ok {
		if role == "manager" || role == "admin" {
			authorized = true
		}
	}

	if verified && authorized {
		err = fire.SetCustomUserClaims(ctx, t.Subject, customClaims)
		L.Fail("claim", "unable to set custom claims", err)

		token, err := fire.VerifyIDToken(ctx, u.IDToken)
		claims := token.Claims

		if manager, ok := claims["manager"]; ok {
			if manager.(bool) {
				L.Good("claim", "manager", "ok")
			}
			return token, nil
		}
		return nil, err
	}
	return nil, err
}

func NewAdminCustomClaims(u *UserCredentials) (*auth.Token, error) {
	// client, err := fire.Auth(ctx)
	// L.Fail("firebs", "unable to get auth client", err)

	customClaims := u.CustomClaims
	t, err := fire.VerifyIDToken(ctx, u.IDToken)
	L.Fail("authv", "unable to verify id token", err)

	authorized := false
	verified := t.UID == u.UID
	is_admin := psql.CheckAdminPrivileges(u.UID)

	if role, ok := t.Claims["role"]; ok {
		if role == "admin" {
			authorized = true
		}
	}

	if verified && authorized && is_admin {
		err = fire.SetCustomUserClaims(ctx, t.Subject, customClaims)
		L.Fail("claim", "unable to set admin claims", err)
		token, err := fire.VerifyIDToken(ctx, u.IDToken)
		L.Fail("claim", "unable to verify id token", err)
		claims := token.Claims

		if admin, ok := claims["admin"]; ok {
			if admin.(bool) {
				L.Good("claim", "admin", "ok")
			}
			return token, nil
		}
		return nil, err
	}
	return nil, err
}

func NewOneTimeClaim(u *UserCredentials) (*auth.Token, error) {
	// client, err := fire.Auth(ctx)
	// L.Fail("firebs", "unable to get auth client", err)

	custom_claims := map[string]interface{}{"admin": true}
	t, err := fire.VerifyIDToken(ctx, u.IDToken)
	L.Fail("authv", "unable to verify id token", err)

	verified := t.UID == u.UID

	if verified {
		L.Info("claim", "one-time-claim", t.Issuer)
		err = fire.SetCustomUserClaims(ctx, t.Subject, custom_claims)
		L.Fail("claim", "unable to set admin claims", err)
		L.Good("claim", "one-time-claim", "ok", err)

		token, err := fire.VerifyIDToken(ctx, u.IDToken)
		L.Fail("claim", "unable to verify id token", err)
		claims := token.Claims

		if admin, ok := claims["admin"]; ok {
			if admin.(bool) {
				L.Good("claim", "admin", "ok")
			}
			return token, nil
		}
		return nil, err
	}

	return nil, err
}
