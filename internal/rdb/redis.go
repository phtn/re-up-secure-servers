package rdb

import (
	"context"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"

	"fast/config"
	"fast/internal/models"
	"fast/pkg/utils"

	"firebase.google.com/go/v4/auth"
)

var (
	r   = utils.Ice(" 𝐑 ", 0)
	Rdb = config.LoadConfig().Rdbs
	L   = utils.NewConsole()
)

func RedisHealth() interface{} {
	start := time.Now()
	ctx := context.Background()
	ping := Rdb.Ping(ctx)
	L.Good(r, "ping", ping, nil)
	elapsed := time.Now().Sub(start) / time.Millisecond
	response := map[string]interface{}{
		"sys":     r,
		"elapsed": elapsed,
		"unit":    "ms",
	}
	return response
}

func StoreUserTokens(ctx context.Context, u *Tokens, expiresIn *time.Duration) error {
	key := "usr::" + u.UID + "::token"
	expiry := time.Now().Add(time.Duration(expiresIn.Seconds()))
	value := map[string]interface{}{
		"id_token":      u.IDToken,
		"refresh_token": u.Refresh,
		"uid":           u.UID,
		"expiry":        expiry,
	}
	err := Rdb.Set(ctx, key, value, 0).Err()
	L.Fail(r, "store-refresh", err)
	if err != nil {
		return err
	}
	return nil
}

func GetUserTokens(uid string) (*Tokens, error) {
	ctx := context.Background()
	key := "usr::" + uid + "::token"
	t, err := Rdb.Get(ctx, key).Result()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			L.Warn(r, "get-tokens key-not-found", err)
			return nil, nil
		} else {
			L.Fail(r, "redis-error", err)
			return nil, err
		}
	}

	var v *Tokens
	err = json.Unmarshal([]byte(t), &v)
	L.Fail(r, "store-unmarshal", err)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func StoreToken(k string, f string, t *auth.Token) {

	value, err := json.Marshal(&t)
	L.Fail(r, "storeToken", err)

	ctx := context.Background()
	err = Rdb.Set(ctx, k, value, 6*time.Hour).Err()
	L.Fail(r, "set store-token", err)
	L.Good(r, "set store-token", k, err)
}

func GetActivation(k string) *models.ActivationResponse {
	ctx := context.Background()

	v, err := Rdb.Get(ctx, k).Result()
	L.Fail(r, "result", err)
	var act *models.ActivationResponse
	err = json.Unmarshal([]byte(v), &act)
	L.Fail(r, "unmarshal", err)

	return act
}

func RetrieveToken(k string) (*auth.Token, error) {
	ctx := context.Background()

	val, err := Rdb.Get(ctx, k).Result()
	L.Fail(r, "retrieve-token", err)

	var token *auth.Token

	err = json.Unmarshal([]byte(val), &token)
	L.Fail(r, "retrieveToken", err)

	L.Good(r, "retrievetoken", "done", err)
	return token, err
}

func StoreVal(key string, h time.Duration, v interface{}) *StoreInfo {
	ctx := context.Background()

	value, err := json.Marshal(v)
	L.Fail(r, "json-marshal", err)

	pipe := Rdb.Pipeline()
	set_cmd := pipe.Set(ctx, key, value, h*time.Minute)
	ttl_cmd := pipe.TTL(ctx, key)

	_, err = pipe.Exec(ctx)
	L.Fail(r, "pipe-exec", "store-value", err)

	err = set_cmd.Err()
	L.Fail(r, "set-cmd", "store-value", err)

	ttl, err := ttl_cmd.Result()
	L.Fail(r, "ttl-cmd", ttl, err)

	L.Good(r, "store-val", key, time.Now().Add(h*time.Hour).Local().UTC().Format("2006-01-02 15:04:05"))
	return &StoreInfo{StoreKey: key, TTL: time.Duration(ttl.Milliseconds())}
}

func RetrieveVal(key string) *RetrieveInfo {
	ctx := context.Background()

	pipe := Rdb.Pipeline()
	get_cmd := pipe.Get(ctx, key)
	ttl_cmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	L.Fail(r, "pipe-exec", "retrieve-value", key, err)

	val, err := get_cmd.Result()
	if val == "" {
		L.Fail(r, "get-cmd", "retrieve-val", key, err)
		return &RetrieveInfo{Exists: false}
	}

	ttl, err := ttl_cmd.Result()
	L.Fail(r, "ttl-cmd", "retrieve-val", "no value", err)

	var v *interface{}

	err = json.Unmarshal([]byte(val), &v)
	L.Fail(r, "json-unmarshal", "retrieve-val", err)

	L.Good(r, "retrieve-val", "all good", err)

	return &RetrieveInfo{Exists: true, Value: v, TTL: time.Duration(ttl.Milliseconds())}
}

func DevSet(key string, v auth.Token) string {
	r := "dev"

	value, err := json.Marshal(v)
	L.Fail(r, "json-marshal", err)

	ctx := context.Background()
	err = Rdb.Set(ctx, key, value, 1*time.Hour).Err()
	L.Fail(r, "rdb-set", err)

	L.Good(r, "rdb-set", key, err)
	return key
}

func DevGet(key string) (interface{}, bool) {
	ctx := context.Background()

	data, err := Rdb.Get(ctx, key).Result()
	errIsNil := isNil(err)
	L.Fail(r, "get", err)

	var v auth.Token
	err = json.Unmarshal([]byte(data), &v)
	L.Fail(r, "get token", "key:", key, err)

	L.Good(r, "rdb-get", key, err)
	return v, errIsNil
}

func HashSet() interface{} {
	start := time.Now()
	f := "hash-set"
	ctx := context.Background()

	descartes := map[string]string{
		"I":   "I",
		"II":  "think",
		"III": "Therefore",
		"IV":  "I",
		"V":   "am",
	}

	for k, v := range descartes {
		err := Rdb.HSet(ctx, f, k, v).Err()
		L.Fail(r, f, err)
	}

	val, err := Rdb.Get(ctx, f).Result()
	L.Fail(r, f, err)
	L.Good(r, f, val, err)

	elapsed := time.Now().Sub(start)
	L.Good(r, f, elapsed, err)

	return elapsed
}

func Int_Token_Set(key string, v interface{}) interface{} {
	rte := utils.Fuchsia("𝕊𝕋𝔸𝔾𝔼𝕟", 0)

	value, err := json.Marshal(v)
	L.Fail(rte, "-json-marshal", err)

	ctx := context.Background()
	err = Rdb.Set(ctx, key, value, 24*5*time.Hour).Err()
	L.Fail(rte, "debug-rdb-set", err)

	L.Good(rte, "rdb token set", key, err)
	return key
}

func Int_Token_Get(key string) (interface{}, bool) {
	ctx := context.Background()
	rte := utils.Dev("𝕊𝕋𝔸𝔾𝔼𝕟", 0)

	data, err := Rdb.Get(ctx, key).Result()
	errIsNil := isNil(err)
	L.Fail(rte, "get", err)

	var v auth.Token
	err = json.Unmarshal([]byte(data), &v)
	L.Fail(rte, "get token", "key:", key, err)

	L.Good(rte, "rdb-get", key, err)
	return v, errIsNil
}

func DebugSet(key string, v interface{}) interface{} {

	value, err := json.Marshal(v)
	L.Fail(dev, "debug-json-marshal", err)

	ctx := context.Background()
	err = Rdb.Set(ctx, key, value, 1*time.Minute).Err()
	L.Fail(dev, "debug-rdb-set", err)

	L.Good(dev, "debug-rdb-set", key, err)
	return key
}
