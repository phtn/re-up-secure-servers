package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	math "math/rand"
	"strings"
)

var (
	div         = Gray("ꔷꔷ", 0)
	indicator   = "●"
	f_indicator = "⨵"
	d_indicator = "◌"
)

func GenerateCode() (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const length = 6

	code := make([]byte, length)
	for i := 0; i < length; i++ {
		// Generate cryptographically secure random number
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		code[i] = letters[n.Int64()]
	}

	return string(code), nil
}

func Struffle(i string) string {
	if len(i) != 6 {
		return "Input must be 6 characters long"
	}
	return string([]byte{i[4], i[1], i[2], i[0], i[3], i[5]})
}

func s() string {
	return strings.ToLower(fmt.Sprintf("%04x", math.Intn(0x10000))[1:])
}

func Guid() string {
	return fmt.Sprintf("%s%s-%s-%s-%s-%s%s%s",
		s(), s(), s(), s(), s(), s(), s(), s())
}

func DStruct[T any](generic interface{}) (T, error) {
	if v, ok := generic.(T); ok {
		return v, nil
	}
	var zeroval T
	return zeroval, fmt.Errorf("type assertion failed")
}
