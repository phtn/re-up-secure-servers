package models

import (
	"fast/config"
	"fast/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

var (
	L = utils.NewConsole()
)

type Role string

const (
	Dev   Role = "developer"
	Adm   Role = "admin"
	Mgr   Role = "manager"
	Acm   Role = "account_manager"
	Agent Role = "agent"
	Suprv Role = "supervisor"
	Undrw Role = "underwriter"
	Dealr Role = "dealer"
)

type Account struct {
	UID         string `json:"uid,omitempty"`
	Name        string `json:"name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
	APIKey      string `json:"api_key,omitempty"`
	Active      bool   `json:"is_active,omitempty"`
	CrTime      string `json:"creation_time,omitempty"`
	AddressId   string `json:"address_id,omitempty"`
	UpdateAt    string `json:"update_time,omitempty"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UID   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role  Role   `json:"role,omitempty"`
}

var (
	jwts = []byte(config.LoadConfig().JwtSecret)
	fire = config.LoadConfig().Fire
	rdb  = config.LoadConfig().Rdbs
)

// func verifyToken(tokenString string) (string, error) {

// 	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwts, nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// 		storeClaims(shield.NewClaimsKey(claims.Email), claims)
// 	}

// 	return "", fmt.Errorf("invalid token")
// }

// func storeClaims(key string, v *CustomClaims) {

// 	value, err := json.Marshal(&v)
// 	L.Fail("json", "marshal-store-claims", err)

// 	ctx := context.Background()
// 	err = rdb.Set(ctx, key, value, 24*7*time.Hour).Err()
// 	L.Fail("redis", "set store-claims", err)

// 	L.Good("redis", "set complete", key, err)
// }

// func GenerateRefreshToken(user Account, key []byte) (string, error) {
// 	claims := jwt.RegisteredClaims{
// 		Subject:   user.UID,
// 		Issuer:    "re-up.ph secure servers",
// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		NotBefore: jwt.NewNumericDate(time.Now()),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
// 		Audience:  []string{"Re-up Secure Servers Clients"},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(key)
// }

// func GenerateAccessToken(user Account) (string, error) {
// 	claims := jwt.RegisteredClaims{
// 		Subject:   user.UID,
// 		Issuer:    "re-up.ph secure servers",
// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		NotBefore: jwt.NewNumericDate(time.Now()),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
// 		Audience:  []string{"Re-up Secure Servers Clients"},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(jwts)
// }
