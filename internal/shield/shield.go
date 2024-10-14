package shield

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fast/pkg/utils"
	"io"
	"os"
	"strconv"
	"strings"
)

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
	utils.ErrLog(r, "shash", err)

	gcm, err := cipher.NewGCM(block)
	utils.ErrLog(r, "gcm", err)

	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	ctext := gcm.Seal(nonce, nonce, value, nil)
	return ctext
}

func Decrypt(c []byte, keyPhrase string) []byte {
	hash := HashIt(keyPhrase)
	block, err := aes.NewCipher([]byte(hash))
	utils.ErrLog("decrypt", "aes", err)

	gcm, err := cipher.NewGCM(block)
	utils.Fatal("decrypt", "gcm", err)
	nsize := gcm.NonceSize()
	nonce, ctext := c[:nsize], c[nsize:]

	original, err := gcm.Open(nil, nonce, ctext, nil)
	utils.Fatal("decrypt", "gcm-open", err)
	utils.Ok("decrypt", "open", strings.Split(string(original), "_")[0])

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

func NewKey(i string) string {
	sep := "--"
	ids := issuerIds()
	idx := utils.RandIdx(48)
	iid := ids[idx]
	key := iid + sep + i + sep + strconv.Itoa(idx)
	return key
}
