package service

import (
	"context"
	"fast/internal/models"
	"fast/pkg/utils"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func CreateCustomToken(uid models.Uid, ctx context.Context, app *firebase.App) string {

	client, err := app.Auth(context.Background())
	if err != nil {
		utils.Err("Error getting auth client %s: %v\n", "[CreateCustomToken]", err)
	}

	token, err := client.CustomToken(context.Background(), uid.UID)
	if err != nil {
		log.Fatalf("Error creating custom token for user %s: %v\n", uid, err)
	}

	utils.Ok("post", "createToken", token)
	return token

}

func VerifyIDToken(ctx context.Context, app *firebase.App, idToken string) *auth.Token {

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("Error verifying ID token: %v\n", err)
	}

	return token

}

func GetUser(ctx context.Context, app *firebase.App, uid models.Uid) *auth.UserRecord {

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting auth client: %v\n", err)
	}

	usr, err := client.GetUser(ctx, uid.UID)
	if err != nil {
		log.Fatalf("Error getting user %s: %v \n", uid.UID, err)
	}

	utils.Ok("get", "getUser", uid.UID)
	return usr

}

func CreateUser(ctx context.Context, client *auth.Client) *auth.UserRecord {

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
		log.Fatalf("error creating user: %v\n", err)
	}

	log.Printf("Successfully created user: %v\n", usr)
	return usr
}
