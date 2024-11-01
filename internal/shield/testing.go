package shield

import (
	"fast/config"
	"fast/pkg/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	L = utils.NewConsole()
	S = config.LoadConfig().JwtSecret
)

func TestEncoding(label string) {
	str := []byte(config.LoadConfig().JwtSecret)
	encoded := EncodeBase64(str)
	L.Info(label, "encoded", encoded)
}

func TestBase64(label string) {
	str := []byte(label)
	encoded := EncodeBase64(str)
	encrypted := Encrypt(str, S)
	L.Info(label, "encoded", encoded)
	L.Info("ecrypted", label, string(encrypted))
	decoded := DecodeBase64(encoded)
	decrypted := Decrypt(encrypted, S)
	L.Info(label, "decoded", string(decoded))
	L.Info("decrypted", label, string(decrypted))
}

func TestJWT(label string) {
	jsonwebtoken := NewJWTKeyFromStr(S)
	encoded := EncodeBase64(jsonwebtoken)
	L.Info(label, "encoded", encoded)
	L.Info(label, "decoded", DecodeBase64(encoded))
}

func TestJWTClaims(label string) string {

	now := time.Now()
	claims := SuperClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "re-up.ph secure servers",
			Subject:   S,
			Audience:  []string{"re-up secure servers test-clients"},
		},
		UID:    NewKey(S),
		Claims: "admin-test-claim",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(NewKey(S)))
	L.Fail("claim", "unable to signed token", err)

	part := strings.Split(tokenString, ".")

	L.Info(label, "part 0", part[0])
	L.Info(label, "encoded", EncodeBase64([]byte(part[0])))
	L.Info(label, "part 1", part[1])
	L.Info(label, "part 1", part[2])

	return tokenString
}
