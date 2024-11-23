package shield

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fast/config"
	"io"
	math "math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type NewAccountToken struct {
	UID   string `json:"uid,omitempty"`
	Email string `json:"email,omitempty"`
}

type SuperClaims struct {
	jwt.RegisteredClaims
	UID    string `json:"uid,omitempty"`
	Claims string `json:"claims,omitempty"`
}

var r = "encrypt"

func HashIt(i string) string {
	input := []byte(i)
	hash := md5.Sum(input)
	return hex.EncodeToString(hash[:])
}

func ShashIt(i string) string {
	input := []byte(i)
	sha := sha256.Sum256(input)
	return hex.EncodeToString(sha[:])
}

func ShashItGood(i string) string {
	input := []byte(i)
	sha := sha512.Sum512(input)
	return hex.EncodeToString(sha[:])
}

func Encrypt(value []byte, keyPhrase string) []byte {
	block, err := aes.NewCipher([]byte(HashIt(keyPhrase)))
	L.Fail(r, "shash", err)

	gcm, err := cipher.NewGCM(block)
	L.Fail(r, "gcm", err)

	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	ctext := gcm.Seal(nonce, nonce, value, nil)
	return ctext
}

func Decrypt(c []byte, keyPhrase string) []byte {
	hash := HashIt(keyPhrase)
	block, err := aes.NewCipher([]byte(hash))
	L.Fail("decrypt", "aes", err)

	gcm, err := cipher.NewGCM(block)
	L.Fail("decrypt", "gcm", err)
	nsize := gcm.NonceSize()
	nonce, ctext := c[:nsize], c[nsize:]

	original, err := gcm.Open(nil, nonce, ctext, nil)
	L.Fail("decrypt", "gcm-open", err)
	L.Good("decrypt", "open", strings.Split(string(original), "_")[0])

	return original
}

func issuerIds() []string {
	issuerId := os.Getenv("RE_UP_ISSUER_ID")
	var ids []string
	for i := 0; i < len(issuerId)-15; i++ {
		ids = append(ids, issuerId[i:i+16])
	}
	return ids
}

func NewClaimsKey(i string) string {
	sep := "--"
	src := uuid.New().String()
	pfx := EncodeBase64([]byte(src))[:16]
	key := pfx + sep + i + sep + "cc"
	return key
}

func NewKey(i string) string {
	sep := "--"
	ids := issuerIds()
	idx := RandIdx(48)
	iid := ids[idx]
	key := iid + sep + i + sep + strconv.Itoa(idx)
	return key
}

func NewJWTSecret() []byte {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	L.Fail("jwt", "create-new", err)
	return b
}

func NewJWTKeyFromStr(secret string) []byte {
	if secret == "" {
		L.Good("jwt", "from-string", "secret must not be an empty string.")
	}
	h := sha256.New()
	h.Write([]byte(secret))
	return h.Sum(nil)
}

func EncodeBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func DecodeBase64(s string) []byte {
	key, err := base64.StdEncoding.DecodeString(s)
	L.Fail("base64", "decode", err)
	return key
}

func NewAccount(u *NewAccountToken) (interface{}, error) {

	now := time.Now()
	claims := SuperClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "re-up.ph secure servers",
			Subject:   u.UID,
			Audience:  []string{"re-up secure servers test-clients"},
		},
		UID:    u.UID,
		Claims: "manager,group-master,active",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(NewKey(S)))
	L.FailR("new acct", "unable to signed token", tokenString, err)

	part := strings.Split(tokenString, ".")

	L.Info("new acct", "token", part[1])
	L.Info("new acct", "key", EncodeBase64([]byte(part[1])))

	result := make(map[string]interface{})
	result["key"] = EncodeBase64([]byte(part[1]))
	result["token"] = part[1]

	return result, nil
}

func RandIdx(n int) int {
	src := math.NewSource(time.Now().UnixNano())
	math.New(src)
	return math.Intn(n)
}

func CreateActivationKey(key_code string) string {
	enc := Encrypt([]byte(key_code), config.LoadConfig().ApiKey)
	eb64 := EncodeBase64(enc)
	return eb64
}
