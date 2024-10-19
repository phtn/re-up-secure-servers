package service

import (
	"context"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/repository"
	"fast/internal/shield"
	"fast/pkg/utils"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

const (
	GET   = "get"
	POST  = "post"
	PATCH = "patch"
)

func eqc(k string, verified bool, t *auth.Token) models.VResult {
	if !verified {
		return models.VResult{
			Verified: verified,
		}
	}
	return models.VResult{
		Key:      k,
		Verified: verified,
		Exp:      int16(t.Expires),
	}
}

func VerifyIdToken(ctx context.Context, fire *firebase.App, v models.VerifyToken) models.VResult {
	var f, r = "verifyIdToken", POST

	utils.Info(v.IDToken[:8], v.Email, v.UID)
	client, err := fire.Auth(context.Background())
	utils.ErrLog(r, f, err)

	k := shield.NewKey(v.Email)
	t, err := client.VerifyIDToken(ctx, v.IDToken)
	utils.ErrLog(r, f, err)

	verified := t.UID == v.UID
	utils.Info("verif", "id_token", verified)

	rdb.StoreToken(k, f, t)
	return eqc(k, verified, t)
}

func VerifyAuthKey(ctx context.Context, fire *firebase.App, v models.VerifyWithAuthKey) models.VResult {
	var r, f = POST, "authk"

	client, err := fire.Auth(context.Background())
	utils.ErrLog(r, f, err)

	k := v.AuthKey
	token, err := rdb.RetrieveToken(k)
	utils.ErrLog(r, f, err)

	verified := false

	if token == nil {
		t, err := client.VerifyIDToken(ctx, v.IDToken)
		utils.ErrLog(r, f, err)
		utils.OkLog(r, f, "verified", err)

		verified = t.UID == v.UID
		return eqc(k, verified, t)
	}

	verified = token.UID == v.UID
	utils.Ok("verif", "id_token", verified)

	return eqc(k, verified, token)
}

func VerifyAdmin(ctx context.Context, fire *firebase.App, v *UserCredentials) bool {
	var r, f = "verif", "admin"

	client, err := fire.Auth(context.Background())
	utils.ErrLog(r, f, err)

	if v.IDToken == "" {
		return false
	}

	t, err := client.VerifyIDToken(ctx, v.IDToken)
	utils.ErrLog(r, f, err)

	verified := t.UID == v.UID
	with_claims := false
	claims := t.Claims
	if custom_claims, ok := claims["manager"]; ok {
		if custom_claims.(bool) {
			utils.Ok("claim", "manager", "ok")
			return verified && ok
		}
		with_claims = ok
		return ok
	}
	if admin_claims, ok := claims["admin"]; ok {
		if admin_claims.(bool) {
			utils.Ok("claim", "admin", "ok")
			return verified && ok
		}
		with_claims = ok
		return ok
	}

	utils.Warn(r, f, verified && with_claims)
	return verified && with_claims
}

func CreateToken(uid models.Uid, ctx context.Context, fire *firebase.App) string {
	var f, r = "createToken", POST

	client, err := fire.Auth(context.Background())
	utils.ErrLog(r, f, err)

	token, err := client.CustomToken(context.Background(), uid.UID)
	utils.ErrLog(r, f, err)

	if len(token) >= 8 {
		utils.Ok(r, f, uid.UID+string(" · "+repository.ClrBt)+token[:16]+string(repository.Reset))
	}
	return token
}

func GetUser(ctx context.Context, fire *firebase.App, uid models.Uid) *auth.UserRecord {
	var f, r = "getUser", POST

	client, err := fire.Auth(context.Background())
	utils.ErrLog(r, f, err)

	usr, err := client.GetUser(ctx, uid.UID)
	utils.ErrLog(r, f, err)

	if usr != nil {
		utils.Ok(POST, f, uid.UID[:8])
	}
	return usr
}

func CreateUser(ctx context.Context, client *auth.Client) *auth.UserRecord {
	var f = "createUser"

	params := (&auth.UserToCreate{}).
		Email("user@example.com").
		EmailVerified(false).
		PhoneNumber("+15555550100").
		Password("secretPassword").
		DisplayName("John Doe").
		PhotoURL("http://www.example.com/12345678/photo.png").
		Disabled(false)

	usr, err := client.CreateUser(ctx, params)
	utils.ErrLog(POST, f, err)

	utils.Ok(POST, f, usr)
	return usr
}
