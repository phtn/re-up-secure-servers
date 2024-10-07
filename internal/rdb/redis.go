package rdb

import (
	"context"
	"encoding/json"

	"fast/config"
	"fast/pkg/utils"

	"firebase.google.com/go/v4/auth"
)

var rdb = config.LoadConfig().Rdb

func StoreToken(k string, t *auth.Token) {

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

}

func RetrieveToken(k string) *auth.Token {

	ctx := context.Background()

	val, err := rdb.Get(ctx, k).Result()
	if err != nil {
		utils.ErrLog("get", "retrieveToken", err)
	}

	var token *auth.Token

	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		utils.ErrLog("json", "retrieveToken", err)
	}

	if err == nil {
		utils.OkLog("get", "retrieveToken", "all good")
	}

	return token

}
