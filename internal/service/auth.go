package service

import (
	"context"
	"fast/config"
	"fast/internal/models"
	"fast/internal/psql"
	"fast/internal/rdb"
	"fast/internal/shield"
	"fast/pkg/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

const (
	GET   = "get"
	POST  = "post"
	PATCH = "patch"
)

var (
	fire = config.LoadConfig().Fire.AuthClient
	L    = utils.NewConsole()
	r    = utils.Ice(" 𝐒 ", 0)
)

func eqc(k string, verified bool, t *auth.Token, is_active bool, group string) (*models.VResult, error) {
	if !verified {
		return &models.VResult{
			IsActive: false,
		}, nil
	}
	return &models.VResult{
		IDToken:   "",
		Expiry:    int16(t.Expires),
		UID:       t.UID,
		IsActive:  is_active,
		GroupCode: group,
	}, nil
}

func VerifyIdToken(ctx context.Context, out models.VerifyToken) (*models.VResult, error) {

	t, err := fire.VerifyIDToken(ctx, out.IDToken)
	L.Fail(r, "verification", err)

	if strings.Contains(err.Error(), "expired") {
		L.Warn(r, "Token Expired", err)
		return nil, err
	}

	is_expired := expiryCheck(err)

	if is_expired {
		token, err := fire.VerifyIDTokenAndCheckRevoked(ctx, out.Refresh)
		L.Fail(r, "refresh-token", err)

		if token.UID == t.UID {
			out = models.VerifyToken{IDToken: out.Refresh}
		}
	}

	// cookie, err := fire.SessionCookie(ctx, out.IDToken, 96*3*time.Hour)
	// L.Fail(r, f, "session-cookie", err)

	verified := t.UID == out.UID
	L.Info(r, "verify id-token", verified)

	exists, user := psql.CheckIfUserExists(t.UID)
	is_active := true

	if verified && !exists {
		phone_number := MockPhone()
		if out.GroupCode == "" {
			neo_user := psql.NewUser(out.Email, out.Email, phone_number, out.UID, "NEO")
			L.Info(r, "neo-user insert", neo_user)
			return eqc("neo", verified, t, is_active, "NEO")
		}
		new_user := psql.NewUser("BrightOne", out.Email, phone_number, out.UID, out.GroupCode)
		L.Info(r, "sign-up", new_user)

		if new_user != "" {
			claim := CustomClaims{
				"agent": "true",
			}

			err := fire.SetCustomUserClaims(ctx, t.UID, claim)
			L.Fail(r, "add-claim agent", err)
			return eqc(t.UID, verified, t, is_active, "NEO")
		}
	}
	return eqc(t.UID, verified, t, is_active, user.GroupCode)
}

func GetUserRecord(ctx context.Context, v *models.VerifyToken) *models.Verified {

	user, err := fire.GetUser(ctx, v.UID)
	L.Fail(r, "get-user-record", err)

	verified := user.UID == v.UID
	L.Info(r, user.UID, verified)

	return &models.Verified{
		Verified:   verified,
		UserRecord: user,
	}
}

func GetUserRecordByUID(ctx context.Context, uid string) (models.Verified, error) {

	user_record, err := fire.GetUser(ctx, uid)
	L.Fail(r, "get-user by-id", err)

	verified := user_record.UID == uid
	L.Info(r, user_record.UID, verified)

	response := models.Verified{
		Verified:   verified,
		UserRecord: user_record,
	}

	return response, nil
}

func VerifyAuthKey(ctx context.Context, v models.VerifyWithAuthKey) (*models.VResult, error) {

	k := v.AuthKey
	token, err := rdb.RetrieveToken(k)
	L.Fail(r, "auth-key", err)

	verified := false

	if token == nil {
		t, err := fire.VerifyIDToken(ctx, v.IDToken)
		L.Fail(r, "verificaton", err)

		verified = t.UID == v.UID
		return eqc(k, verified, t, false, "")
	}

	verified = token.UID == v.UID
	L.Good(r, token.UID, verified)

	return eqc(k, verified, token, false, "")
}

func VerifyAdmin(ctx context.Context, v *UserCredentials) bool {

	if v.IDToken == "" {
		return false
	}
	t, err := fire.VerifyIDToken(ctx, v.IDToken)
	L.Fail(r, "admin", err)

	verified := t.UID == v.UID
	with_claims := false
	claims := t.Claims
	if custom_claims, ok := claims["manager"]; ok {
		if custom_claims.(bool) {
			L.Good(r, "claims-manager", "ok")
			return verified && ok
		}
		with_claims = ok
		return ok
	}
	if admin_claims, ok := claims["admin"]; ok {
		if admin_claims.(bool) {
			L.Good(r, "admin", "ok")
			return verified && ok
		}
		with_claims = ok
		return ok
	}

	return verified && with_claims
}

func GetUserInfo(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := fire.GetUser(ctx, uid)
	L.Fail(r, "get-user", "info", err)
	timestamp := user.TokensValidAfterMillis / 1000
	L.Info(r, "validity timestamp", timestamp, "seconds")
	return user, nil
}

func TokenVerification(ctx context.Context, v rdb.UserTokens) (*models.TokenResponse, error) {
	t, err := fire.VerifyIDToken(ctx, v.IDToken)
	L.Fail(r, "service token-verification", err)
	is_expired := expiryCheck(err)

	var response *models.TokenResponse

	if is_expired {
		token, err := fire.VerifyIDTokenAndCheckRevoked(ctx, v.Refresh)
		L.Fail(r, "refresh-token", v.Refresh, err)

		if err != nil {
			L.Fail(r, "checked-revoked", v.Refresh, err)
			return nil, err
		}

		response = &models.TokenResponse{
			Token:    token,
			Verified: token.UID == v.UID,
		}
		return response, nil
	}

	response = &models.TokenResponse{
		Token:    t,
		Verified: t.UID == v.UID,
	}
	return response, nil
}

func UnlockWithKey(key_code string) *models.ActivationResponse {
	key := shield.EncodeActivationKey(key_code)
	a := rdb.GetActivation(key)

	L.Good(r, "unlock-with-key", r)
	return a
}

func expiryCheck(err error) bool {
	expired := false
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			expired = true
		}
		return expired
	}
	return expired
}

// func CreateAgentCode(v models.VerifyToken) AgentCodeResponse {
// 	key := shield.NewKey(v.Email)
// 	url := "https://fastinsure.tech/new/agent/code?key=" + key
// 	rdb.StoreVal(key, 48, url)
// 	L.Info("create  ", "agent", "code", url)
// 	response := AgentCodeResponse{Key: key, URL: url}
// 	return response
// }

func NewToken(uid models.Uid, ctx context.Context, fire *firebase.App) string {

	client, err := fire.Auth(context.Background())
	L.Fail(r, "new-token", err)

	token, err := client.CustomToken(context.Background(), uid.UID)
	L.Fail(r, "custom-token", err)

	if len(token) >= 8 {
		L.Good(r, "good token", token[:16])
	}
	return token
}

func GetUser(ctx context.Context, uid models.Uid) *auth.UserRecord {

	usr, err := fire.GetUser(ctx, uid.UID)
	L.Fail(r, "get-user", err)

	if usr != nil {
		L.Good(r, "Good", uid.UID[:8])
	}
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
	L.Fail(r, "create-user", err)
	return usr
}

func MockPhone() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed the random number generator
	random := 100000 + r.Intn(900000)

	random_str := strconv.Itoa(random)
	phone := "+63915" + random_str + "0"
	return phone
}
