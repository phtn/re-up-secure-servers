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

type Middleware func(http.HandlerFunc) http.HandlerFunc

const (
	Root  = "/v1"
	Auth  = Root + "/auth"
	Admin = Root + "/admin"
)

const (
	AuthPath          = Auth
	GetUserPath       = Auth + "/getUser"
	CreateTokenPath   = Auth + "/createToken"
	VerifyIdTokenPath = Auth + "/verifyIdToken"
	VerifyAuthKeyPath = Auth + "/verifyAuthKey"
	CustomClaimsPath  = Auth + "/createCustomClaims"
	DevSetPath        = Auth + "/devSet"
	DevGetPath        = Auth + "/devGet"
	// ADMIN
	AdminPath       = Admin
	AdminClaimsPath = Admin + "/adminClaims"
)

var (
	fire   = config.LoadConfig().Fire
	secret = []byte(config.LoadConfig().JwtSecret)
)

func DbCheck(w http.ResponseWriter, r *http.Request) {
	psql := models.PsqlHealth()
	redis := rdb.RedisHealth()
	utils.JsonResponse(w, map[string]interface{}{
		"psql": psql, "redis": redis,
	})
}

func VerifyIdToken(w http.ResponseWriter, r *http.Request) {

	utils.PostMethodOnly(w, r)

	// token := jwt.New(jwt.SigningMethodHS256)
	// singed, err := token.SignedString(secret)
	// utils.HttpError(w, "Could not generate token", err, http.StatusInternalServerError)
	// utils.ErrLog("req", "verifyIdToken", err)

	// utils.Ok("auth", "signed", singed)

	body, err := io.ReadAll(r.Body)
	utils.ErrLog("req", "verifyIdToken", err)

	var v models.VerifyToken
	err = json.Unmarshal(body, &v)
	utils.ErrLog("json", "verifyIdToken", err)

	utils.JsonResponse(w, service.VerifyIdToken(r.Context(), fire, v))
}

func VerifyAuthKey(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body-params", err, http.StatusBadRequest)

	var v models.VerifyWithAuthKey
	err = json.Unmarshal(body, &v)
	utils.HttpError(w, "json-body", err, http.StatusBadRequest)

	utils.JsonResponse(w, service.VerifyAuthKey(r.Context(), fire, v))
}

func CreateCustomClaims(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var u *service.UserCredentials
	err = json.Unmarshal(body, &u)
	utils.HttpError(w, "Unable to read json", err, http.StatusBadRequest)

	data, err := service.NewCustomClaims(u)
	utils.HttpError(w, "Unable to read json", err, http.StatusServiceUnavailable)

	utils.JsonResponse(w, data)
}

func CheckAdminAuthority(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var v *service.UserCredentials
	err = json.Unmarshal(body, &v)
	if err != nil {
		utils.Info("handl", "admin-authority", v)
		utils.ErrLog("handl", "admin-authority", err)
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	verified := service.VerifyAdmin(r.Context(), fire, v)
	utils.HttpError(w, "unable to verify credentials", err, http.StatusServiceUnavailable)

	if !verified {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	data := map[string]interface{}{
		"verified":  verified,
		"timestamp": time.Now().Unix(),
	}

	utils.JsonResponse(w, data)
}

func CreateAdminClaims(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var u *service.UserCredentials
	err = json.Unmarshal(body, &u)
	utils.HttpError(w, "Unable to read json", err, http.StatusBadRequest)

	data, err := service.NewAdminCustomClaims(u)
	utils.HttpError(w, "Unable to read json", err, http.StatusServiceUnavailable)

	utils.JsonResponse(w, data)
}

func CreateOneTimeClaim(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	utils.HttpError(w, "body params invalid", err, http.StatusBadRequest)

	var u *service.UserCredentials
	err = json.Unmarshal(body, &u)
	utils.HttpError(w, "Unable to read json", err, http.StatusBadRequest)

	data, err := service.NewOneTimeClaim(u)
	utils.HttpError(w, "Unable to read json", err, http.StatusServiceUnavailable)

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
