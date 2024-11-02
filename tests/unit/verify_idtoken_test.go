package tests

import (
	"testing"
)

func TestVerifyIdToken(t *testing.T) {

	t.Run("test", func(t *testing.T) {
		got := "account"

		if got != "account" {
			t.Error("pink")
		}

	})

}
