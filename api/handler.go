package api

import (
	"encoding/json"
	"fast/config"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/service"
	"fast/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/v4/auth"
)

const (
	bep = "/v1/auth"
)

const (
	AuthRootPath      = bep
	GetUserPath       = bep + "/getUser"
	CreateTokenPath   = bep + "/createToken"
	VerifyIdTokenPath = bep + "/verifyIdToken"
	VerifyAuthKeyPath = bep + "/verifyAuthKey"
	DevSetPath        = bep + "/devSet"
	DevGetPath        = bep + "/devGet"
)

var fire = config.LoadConfig().Fire

func VerifyIdToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("req", "verifyIdToken", err)

	var v models.VerifyToken
	err = json.Unmarshal(body, &v)
	utils.ErrLog("json", "verifyIdToken", err)

	utils.JsonResponse(w, service.VerifyIdToken(r.Context(), fire, v))
}

func DevSet(w http.ResponseWriter, r *http.Request) {
	dev := "dev"

	body, err := io.ReadAll(r.Body)
	utils.ErrLog(dev, "req-body", err)

	var v auth.Token
	err = json.Unmarshal(body, &v)
	utils.ErrLog(dev, "json-unmarshal", err)

	utils.Ok(dev, "rbd-set", time.Nanosecond)
	utils.JsonResponse(w, rdb.DevSet(devKey, v))
}

func VerifyAuthKey(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("req", "verifyAuthKey", err)

	var v models.VerifyWithAuthKey
	err = json.Unmarshal(body, &v)
	utils.ErrLog("json", "verifyAuthKey", err)

	utils.JsonResponse(w, service.VerifyAuthKey(r.Context(), fire, v))
}

func RDBC(w http.ResponseWriter, r *http.Request) {
	data := rdb.RDBC()
	utils.JsonResponse(w, data)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uid models.Uid
	err = json.Unmarshal(body, &uid)
	if err != nil {
		log.Fatal(err)
	}

	var response = service.CreateToken(uid, r.Context(), fire)
	utils.JsonResponse(w, response)
	log.Println(response)

}

func GetUser(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uid models.Uid
	err = json.Unmarshal(body, &uid)

	if err != nil {
		log.Fatal(err)
	}

	var response = service.GetUser(r.Context(), fire, uid)
	utils.JsonResponse(w, response)
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "retrieveToken", err)

	var k models.AuthKey
	err = json.Unmarshal(body, &k)
	utils.ErrLog("json", "retrieveToken", err)

	result, err := rdb.RetrieveToken(k.FastAuthKey)
	utils.ErrLog("redis", "retrieveToken", err)

	utils.JsonResponse(w, result)
}

func StoreVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "storeVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.ErrLog("json", "storeVal", err)

	utils.JsonResponse(w, rdb.StoreVal(kv.Key, kv.Value))
}

func RetrieveVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("redis", "retrieveVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.ErrLog("json", "retrieveVal", err)

	utils.JsonResponse(w, rdb.RetrieveVal(kv.Key))
}

func XXX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kamusta ka?")
}

var pass = os.Getenv("RDB_PASS")[:6]
var devKey = "__DEV__" + pass

func DevGet(w http.ResponseWriter, r *http.Request) {
	dev := "dev"

	_, err := io.ReadAll(r.Body)
	utils.ErrLog(dev, "req-body", err)

	utils.Ok(dev, "rbd-get", time.Nanosecond)
	utils.JsonResponse(w, rdb.DevGet(devKey))
}
