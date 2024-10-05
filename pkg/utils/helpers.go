package utils

import (
	"encoding/json"
	"fast/internal/repository"
	"log"
	"net/http"
)

func Err(msg string, sub string, err error) {
	log.Fatalf(msg, "%s: %v\n", sub, err)
}

func ErrHandler(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return err
}

func Ok(r string, f string, p interface{}) {

	log.Printf(repository.ColorResp+repository.Success+repository.ColorDark+" ৷ "+repository.ColorCode+r+repository.ColorDark+" ৷ "+repository.ColorReset+f+repository.ColorLogStart+": %s\n", p)
}

func JsonResponse(w http.ResponseWriter, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	ErrHandler(w, err)

	_, err = w.Write(jsonData)
	ErrHandler(w, err)
}
