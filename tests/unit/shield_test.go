package tests

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashIt(t *testing.T) {
	tests := []struct {
		description string
		expected    string
	}{
		{
			description: "hashed the input",
			expected:    "b5c0b187fe309af0f4d35982fd961d7e",
		},
	}

	for _, test := range tests {

		i := "love"
		input := []byte(i)
		hash := md5.Sum(input)
		result := hex.EncodeToString(hash[:])

		assert.Exactlyf(t, test.expected, result, test.description)
	}
}
