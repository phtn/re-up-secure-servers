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
)

var app = config.LoadConfig().App

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

	var response = service.CreateToken(uid, r.Context(), app)
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

	var response = service.GetUser(r.Context(), app, uid)
	utils.JsonResponse(w, response)

}

func VerifyIDToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var idToken models.IdToken
	err = json.Unmarshal(body, &idToken)
	if err != nil {
		utils.ErrLog("json", "verifyIdToken", err)
	}

	utils.JsonResponse(w, service.VerifyIDToken(r.Context(), app, &idToken))
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.CheckErrLog("redis", "retrieveToken", err)

	var k models.AuthKey
	err = json.Unmarshal(body, &k)
	utils.CheckErrLog("json", "retrieveToken", err)

	utils.JsonResponse(w, rdb.RetrieveToken(k.FastAuthKey))
}

func StoreVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.CheckErrLog("redis", "storeVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.CheckErrLog("json", "storeVal", err)

	utils.JsonResponse(w, rdb.StoreVal(kv.Key, kv.Value))
}

func RetrieveVal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	utils.CheckErrLog("redis", "retrieveVal", err)

	var kv models.KV
	err = json.Unmarshal(body, &kv)
	utils.CheckErrLog("json", "retrieveVal", err)

	utils.JsonResponse(w, rdb.RetrieveVal(kv.Key))
}

func XXX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kamusta ka?")
}
