package rdb

import (
	"context"
	"encoding/json"

	"fast/config"
	"fast/pkg/utils"

	"firebase.google.com/go/v4/auth"
)

var rdb = config.LoadConfig().Rdb

func StoreToken(k string, t *auth.Token) string {

	ctx := context.Background()

	jsonData, err := json.Marshal(t)
	if err != nil {
		utils.ErrLog("json", "storeToken", err)
	}

	err = rdb.Set(ctx, k, jsonData, 0).Err()
	if err != nil {
		utils.ErrLog("set", "storeToken", err)
	}

	if err == nil {
		utils.OkLog("set", "storeToken", "all good")
	}
	return k
}

func RetrieveToken(k string) *auth.Token {

	ctx := context.Background()

	val, err := rdb.Get(ctx, k).Result()
	utils.CheckErrLog("get", "retrieveToken", err)

	var token *auth.Token

	err = json.Unmarshal([]byte(val), &token)
	utils.CheckErrLog("json", "retrieveToken", err)

	if err == nil {
		utils.OkLog("get", "retrieveToken", "all good")
	}

	return token
}

func StoreVal(key string, v interface{}) string {

	ctx := context.Background()

	value, err := json.Marshal(v)
	utils.CheckErrLog("json", "storeToken", err)

	err = rdb.Set(ctx, key, value, 0).Err()
	utils.CheckErrLog("set", "storeToken", err)

	if err == nil {
		utils.OkLog("set", "storeToken", "all good")
	}
	return key
}

func RetrieveVal(key string) interface{} {

	ctx := context.Background()

	val, err := rdb.Get(ctx, key).Result()
	utils.CheckErrLog("get", "retrieveToken", err)

	var v *interface{}

	err = json.Unmarshal([]byte(val), &v)
	utils.CheckErrLog("json", "retrieveToken", err)

	if err == nil {
		utils.OkLog("redis", "retrieveVal", "all good")
	}

	return v
}
