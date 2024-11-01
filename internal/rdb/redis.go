package rdb

import (
	"context"
	"encoding/json"
	"time"

	"fast/config"
	"fast/pkg/utils"

	"firebase.google.com/go/v4/auth"
)

type RetrieveInfo struct {
	Exists bool          `json:"exists"`
	Value  interface{}   `json:"value,omitempty"`
	TTL    time.Duration `json:"ttl,omitempty"`
}
type StoreInfo struct {
	StoreKey string        `json:"store_key,omitempty"`
	TTL      time.Duration `json:"ttl,omitempty"`
}

var (
	r   = "redis"
	rdb = config.LoadConfig().Rdbs
	L   = utils.NewConsole()
)

func RedisHealth() interface{} {
	start := time.Now()
	ctx := context.Background()
	ping := rdb.Ping(ctx)
	L.Good(r, "ping", ping, nil)
	elapsed := time.Now().Sub(start) / time.Millisecond
	response := map[string]interface{}{
		"sys":     r,
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

func StoreVal(key string, h time.Duration, v interface{}) *StoreInfo {
	ctx := context.Background()

	value, err := json.Marshal(v)
	L.Fail(r, "json-marshal", err)

	pipe := rdb.Pipeline()
	set_cmd := pipe.Set(ctx, key, value, h*time.Hour)
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

	pipe := rdb.Pipeline()
	get_cmd := pipe.Get(ctx, key)
	ttl_cmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	L.Fail(r, "pipe-exec", "retrieve-value", err)

	val, err := get_cmd.Result()
	if val == "" {
		L.Fail(r, "get-cmd", "retrieve-val", err)
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
		L.Fail("rdb", f, err)
	}

	val, err := rdb.Get(ctx, f).Result()
	L.Fail("rdb", f, err)
	L.Good("rdb", f, val, err)

	elapsed := time.Now().Sub(start)
	L.Good("rdb", f, elapsed, err)

	return elapsed
}
