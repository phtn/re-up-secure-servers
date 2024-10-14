package rdb

import (
	"context"
	"encoding/json"
	"time"

	"fast/config"
	"fast/pkg/utils"

	"firebase.google.com/go/v4/auth"
)

var rdb = config.LoadConfig().Rdbs

func RDBC() interface{} {
	start := time.Now()
	r, f := "rdb", "conn"

	ctx := context.Background()

	ping := rdb.Ping(ctx)
	utils.Ok(r, f, ping)

	return time.Now().Sub(start)
}

func StoreToken(k string, f string, t *auth.Token) {

	value, err := json.Marshal(&t)
	utils.ErrLog("json", "storeToken", err)

	ctx := context.Background()
	err = rdb.Set(ctx, k, value, 1*time.Hour).Err()
	utils.ErrLog("set", "storeToken", err)

	utils.OkLog("rdb", "storeToken", k, err)
}

func RetrieveToken(k string) (*auth.Token, error) {
	ctx := context.Background()

	val, err := rdb.Get(ctx, k).Result()
	utils.ErrLog("get", "retrieveToken", err)

	var token *auth.Token

	err = json.Unmarshal([]byte(val), &token)
	utils.ErrLog("json", "retrieveToken", err)

	utils.OkLog("get", "retrieveToken", "done", err)
	return token, err
}

func StoreVal(key string, v interface{}) string {
	ctx := context.Background()

	value, err := json.Marshal(v)
	utils.ErrLog("json", "storeToken", err)

	err = rdb.Set(ctx, key, value, 1*time.Hour).Err()
	utils.ErrLog("set", "storeToken", err)

	utils.OkLog("set", "storeToken", "all good", err)
	return key
}

func RetrieveVal(key string) interface{} {
	ctx := context.Background()

	val, err := rdb.Get(ctx, key).Result()
	utils.ErrLog("get", "retrieveToken", err)

	var v *interface{}

	err = json.Unmarshal([]byte(val), &v)
	utils.ErrLog("json", "retrieveToken", err)

	utils.OkLog("redis", "retrieveVal", "all good", err)
	return v
}

func DevSet(key string, v auth.Token) string {
	r := "dev"

	value, err := json.Marshal(v)
	utils.ErrLog(r, "json-marshal", err)

	ctx := context.Background()
	err = rdb.Set(ctx, key, value, 1*time.Hour).Err()
	utils.ErrLog(r, "rdb-set", err)

	utils.OkLog(r, "rdb-set", key, err)
	return key
}

func DevGet(key string) interface{} {
	ctx := context.Background()

	data, err := rdb.Get(ctx, key).Result()
	utils.ErrLog("dev", "get", err)

	var v auth.Token
	err = json.Unmarshal([]byte(data), &v)
	utils.ErrLog("json", "retrieveToken", err)

	utils.OkLog("dev", "rdb-get", key, err)
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
	utils.ErrLog("rdb", f, err)
	utils.OkLog("rdb", f, val, err)

	elapsed := time.Now().Sub(start)
	utils.OkLog("rdb", f, elapsed, err)

	return elapsed
}
