package service

import (
	"context"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/repository"
	"fast/pkg/utils"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

const (
	GET   = "get"
	POST  = "post"
	PATCH = "patch"
)

func CreateToken(uid models.Uid, ctx context.Context, app *firebase.App) string {
	var f = "createToken"

	client, err := app.Auth(context.Background())
	if err != nil {
		utils.ErrLog(POST, f, err)
	}

	token, err := client.CustomToken(context.Background(), uid.UID)
	if err != nil {
		utils.ErrLog(POST, f, err)
	}

	if len(token) >= 8 {
		utils.OkLog(POST, f, uid.UID+string(" · "+repository.Bright)+token[:16]+string(repository.Reset))
	}
	return token

}

func VerifyIDToken(ctx context.Context, app *firebase.App, idToken models.IdToken) string {
	var f = "verifyIdToken"

	client, err := app.Auth(context.Background())
	if err != nil {
		utils.ErrLog(POST, f, err)
	}

	k := utils.Guid()
	t, err := client.VerifyIDToken(ctx, idToken.Token)
	if err != nil {
		utils.ErrLog(POST, f, err)
		rdb.StoreToken(k, t)
		utils.OkLog(POST, f, repository.Bright+"token stored"+string(repository.Reset))
	}

	return k

}

func GetUser(ctx context.Context, app *firebase.App, uid models.Uid) *auth.UserRecord {
	var f = "getUser"

	client, err := app.Auth(context.Background())
	if err != nil {
		utils.ErrLog(POST, f, err)
	}

	usr, err := client.GetUser(ctx, uid.UID)
	if err != nil {
		utils.NilLog(POST, f, err)
	}

	if usr != nil {
		utils.OkLog(POST, f, uid.UID[:8])
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
	if err != nil {
		utils.ErrLog(POST, f, err)
	}

	utils.OkLog(POST, f, usr)
	return usr
}
