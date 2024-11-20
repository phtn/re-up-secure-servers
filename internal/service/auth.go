package service

import (
	"context"
	"fast/config"
	"fast/internal/models"
	"fast/internal/psql"
	"fast/internal/rdb"
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
	fire = config.LoadConfig().Fire.Auth
	L    = utils.NewConsole()
)

func eqc(k string, verified bool, t *auth.Token, is_active bool, cookie string) models.VResult {
	if !verified {
		return models.VResult{
			Verified: verified,
		}
	}
	return models.VResult{
		Key:      k,
		Verified: verified,
		Expiry:   int16(t.Expires),
		IsActive: is_active,
		Cookie:   cookie,
	}
}

func VerifyIdToken(ctx context.Context, out models.VerifyToken) models.VResult {
	var r, f = "v-idToken", "service: VerifyIdToken"

	L.Info(out.IDToken[:8], out.Email, out.UID)

	t, err := fire.VerifyIDToken(ctx, out.IDToken)
	is_expired := expiryCheck(err)
	L.Fail(r, f, "is_expired", is_expired, err)

	if is_expired {
		token, err := fire.VerifyIDTokenAndCheckRevoked(ctx, out.IDToken)
		L.Fail("refresh-token", "service: VerifyIdToken", err)

		if token.UID == t.UID {
			out = models.VerifyToken{IDToken: out.Refresh}
		}
	}

	// cookie, err := fire.SessionCookie(ctx, out.IDToken, 96*3*time.Hour)
	// L.Fail(r, f, "session-cookie", err)

	verified := t.UID == out.UID
	L.Info("verify", "id_token", verified)

	exists := psql.CheckIfUserExists(t.UID)
	is_active := true

	if verified && !exists {
		phone_number := mock_phone()
		if out.GroupCode == "" {
			neo_user := psql.NewUser("Neo", out.Email, phone_number, out.UID, "NEO")
			L.Info("neo-user-insert", "service: VerifyIdToken", "uid", neo_user)
			return eqc("neo", verified, t, is_active, "cookie")
		}
		new_user := psql.NewUser("BrightOne", out.Email, phone_number, out.UID, out.GroupCode)
		L.Info("new-user", "sign-up", new_user)

		if new_user != "" {
			claim := CustomClaims{
				"agent": true,
			}
			token, err := AddCustomClaim(out.IDToken, out.UID, claim)
			L.Fail("add-claim", "agent", err)
			L.Good("add-claim", "agent", "success", err)
			return eqc(token.UID, verified, token, is_active, "cookie")
		}
	}
	return eqc(t.UID, verified, t, is_active, "cookie")
}

func GetUserRecord(ctx context.Context, v *models.VerifyToken) *models.Verified {
	var r, f = "id-token", "verified"

	user, err := fire.GetUser(ctx, v.UID)
	L.Fail(r, f, err)

	verified := user.UID == v.UID
	L.Info(r, f, verified)

	return &models.Verified{
		Verified:   verified,
		UserRecord: user,
	}
}

func VerifyAuthKey(ctx context.Context, v models.VerifyWithAuthKey) models.VResult {
	var r, f = POST, "authky"

	k := v.AuthKey
	token, err := rdb.RetrieveToken(k)
	L.Fail(r, f, err)

	verified := false

	if token == nil {
		t, err := fire.VerifyIDToken(ctx, v.IDToken)
		L.Fail(r, f, err)
		L.Good(r, f, "verified", err)

		verified = t.UID == v.UID
		return eqc(k, verified, t, false, "")
	}

	verified = token.UID == v.UID
	L.Good("verify", "id_token", verified)

	return eqc(k, verified, token, false, "")
}

func VerifyAdmin(ctx context.Context, v *UserCredentials) bool {
	var r, f = "verify", "admin"

	// client, err := fire.Auth(context.Background())
	// L.Fail(r, f, err)

	if v.IDToken == "" {
		return false
	}
	t, err := fire.VerifyIDToken(ctx, v.IDToken)
	L.Fail(r, f, err)

	verified := t.UID == v.UID
	with_claims := false
	claims := t.Claims
	if custom_claims, ok := claims["manager"]; ok {
		if custom_claims.(bool) {
			L.Good("claims", "manager", "ok")
			return verified && ok
		}
		with_claims = ok
		return ok
	}
	if admin_claims, ok := claims["admin"]; ok {
		if admin_claims.(bool) {
			L.Good("claims", "admin", "ok")
			return verified && ok
		}
		with_claims = ok
		return ok
	}

	L.Warn(r, f, verified && with_claims)
	return verified && with_claims
}

func GetUserInfo(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := fire.GetUser(ctx, uid)
	L.FailR("get-user", "service", "firebase", err)
	timestamp := user.TokensValidAfterMillis / 1000
	L.Info("token-validity", "timestamp", timestamp, "seconds")
	return user, nil
}

func TokenVerification(ctx context.Context, v models.UserRefresh) (*models.TokenResponse, error) {
	t, err := fire.VerifyIDToken(ctx, v.IDToken)
	L.Fail("fire.VerifyIDToken", "service: TokenVerification", err)
	is_expired := expiryCheck(err)

	var response *models.TokenResponse

	if is_expired {
		token, err := fire.VerifyIDTokenAndCheckRevoked(ctx, v.Refresh)
		L.Fail("refresh-token", "service: TokenVerification", v.Refresh, err)

		if err != nil {
			L.Fail("checked-revoked", v.Refresh, err)
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
	var r, f = "new token", "NewToken"

	client, err := fire.Auth(context.Background())
	L.Fail(r, f, err)

	token, err := client.CustomToken(context.Background(), uid.UID)
	L.Fail(r, f, err)

	if len(token) >= 8 {
		L.Good(r, f, token[:16])
	}
	return token
}

func GetUser(ctx context.Context, uid models.Uid) *auth.UserRecord {
	var f, r = "getUser", POST

	// client, err := fire.Auth(context.Background())
	// L.Fail(r, f, err)

	usr, err := fire.GetUser(ctx, uid.UID)
	L.Fail(r, f, err)

	if usr != nil {
		L.Good(POST, f, uid.UID[:8])
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
	L.Fail(POST, f, err)

	L.Good(POST, f, usr)
	return usr
}

func mock_phone() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed the random number generator
	random := 100000 + r.Intn(900000)

	random_str := strconv.Itoa(random)
	phone := "+63915" + random_str + "0"
	return phone
}
