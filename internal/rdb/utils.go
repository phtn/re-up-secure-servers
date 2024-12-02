package rdb

import (
	"errors"
	"fast/pkg/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserTokens struct {
	IDToken string `json:"id_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	UID     string `json:"uid"`
}

type Tokens struct {
	IDToken string    `json:"id_token,omitempty"`
	Refresh string    `json:"refresh_token,omitempty"`
	Expiry  time.Time `json:"expires_in,omitempty"`
	UID     string    `json:"uid,omitempty"`
}

type TS struct {
	Field  string    `json:"field,omitempty"`
	Value  string    `json:"value,omitempty"`
	Expiry time.Time `json:"expires_in,omitempty"`
}

type RetrieveInfo struct {
	Exists bool          `json:"exists"`
	Value  interface{}   `json:"value,omitempty"`
	TTL    time.Duration `json:"ttl,omitempty"`
}
type StoreInfo struct {
	StoreKey string        `json:"store_key,omitempty"`
	TTL      time.Duration `json:"ttl,omitempty"`
}

type RedisError struct {
	Code    int
	Message string
}

func (e *RedisError) Error() string {
	return fmt.Sprintf("Redis error [%d]: %s", e.Code, e.Message)
}

func isNil(e error) bool {
	if errors.Is(e, redis.Nil) {
		L.Fail(r, "key-not-found", e)
		return true
	}
	return false
}

var (
	dev = utils.Dev("Δ", 0)
)
