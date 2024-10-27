package rdb

import (
	"context"
	"encoding/json"
	"time"

	"fast/config"
	"fast/pkg/utils"

	"firebase.google.com/go/v4/auth"
)

var (
	rdb = config.LoadConfig().Rdbs
	L   = utils.NewConsole()
)

func RedisHealth() interface{} {
	start := time.Now()
	ctx := context.Background()
	ping := rdb.Ping(ctx)
	L.Good("redis", "ping", ping, nil)
	elapsed := time.Now().Sub(start) / time.Millisecond
	response := map[string]interface{}{
		"sys":     "redis",
		"elapsed": elapsed,
		"unit":    "ms",
	}
	return response
}

func StoreToken(k string, f string, t *auth.Token) {

	value, err := json.Marshal(&t)
	L.Fail("json", "storeToken", err)

	ctx := context.Background()
	err = rdb.Set(ctx, k, value, 1*time.Hour).Err()
	L.Fail("set", "storeToken", err)
	L.Good("rdb", "storeToken", k, err)
}

func RetrieveToken(k string) (*auth.Token, error) {
	ctx := context.Background()

	val, err := rdb.Get(ctx, k).Result()
	L.Fail("get", "retrieveToken", err)

	var token *auth.Token

	err = json.Unmarshal([]byte(val), &token)
	L.Fail("json", "retrieveToken", err)

	L.Good("get", "retrieveToken", "done", err)
	return token, err
}

func StoreVal(key string, h time.Duration, v interface{}) string {
	ctx := context.Background()

	value, err := json.Marshal(v)
	L.Fail("json", "marshal", err)

	err = rdb.Set(ctx, key, value, h*time.Hour).Err()
	L.Fail("redis", "set-failed", key, err)

	L.Good("redis", "set-key", key, time.Now().Add(h*time.Hour).Local().UTC().Format("2006-01-02 15:04:05"))
	return key
}

func RetrieveVal(key string) interface{} {
	ctx := context.Background()

	val, err := rdb.Get(ctx, key).Result()
	L.Fail("get", "retrieveToken", err)

	var v *interface{}

	err = json.Unmarshal([]byte(val), &v)
	L.Fail("json", "retrieveToken", err)

	L.Good("redis", "retrieveVal", "all good", err)
	return v
}

func DevSet(key string, v auth.Token) string {
	r := "dev"

	value, err := json.Marshal(v)
	L.Fail(r, "json-marshal", err)

	ctx := context.Background()
	err = rdb.Set(ctx, key, value, 1*time.Hour).Err()
	L.Fail(r, "rdb-set", err)

	L.Good(r, "rdb-set", key, err)
	return key
}

func DevGet(key string) interface{} {
	ctx := context.Background()

	data, err := rdb.Get(ctx, key).Result()
	L.Fail("dev", "get", err)

	var v auth.Token
	err = json.Unmarshal([]byte(data), &v)
	L.Fail("json", "retrieveToken", err)

	L.Good("dev", "rdb-get", key, err)
	return v
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
		err := rdb.HSet(ctx, f, k, v).Err()
		utils.ErrLog("rdb", f, err)
	}

	val, err := rdb.Get(ctx, f).Result()
	L.Fail("rdb", f, err)
	L.Good("rdb", f, val, err)

	elapsed := time.Now().Sub(start)
	L.Good("rdb", f, elapsed, err)

	return elapsed
}
