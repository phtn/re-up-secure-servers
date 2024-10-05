package api

import (
	"encoding/json"
	"fast/config"
	"fast/internal/models"
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

	var response = service.CreateCustomToken(uid, r.Context(), app)
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

func XXX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kamusta ka?")
}
