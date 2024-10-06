package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	Start   = "\033[38;5;60m"
	Success = "\033[38;5;150m"
	Warn    = "\033[38;5;13m"
	Dark    = "\033[38;5;235m"
	Code    = "\033[38;5;153m"
	Error   = "\033[38;5;216m"
	Bright  = "\033[38;5;229m"
	Reset   = "\033[0m"
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

func OkLog(r string, f string, p interface{}) {
	log.Printf(Success+"success"+Dark+" ৷ "+Code+r+Dark+" ৷ "+Reset+f+Start+": %s\n", p)
}

func ErrLog(r string, f string, e error) {
	log.Printf(Warn+"failed"+Dark+"  ৷ "+Code+r+Dark+" ৷ "+Reset+f+Start+": %v\n", e)
}
func NilLog(r string, f string, e error) {
	log.Printf(Warn+"NULL"+Dark+" ·· ৷ "+Code+r+Dark+" ৷ "+Reset+f+Start+": %v\n", e)
}

func JsonResponse(w http.ResponseWriter, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	ErrHandler(w, err)

	_, err = w.Write(jsonData)
	ErrHandler(w, err)
}
