package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	// COLORS
	Start = "\033[38;5;60m"
	ClrOk = "\033[38;5;150m"
	ClrWn = "\033[38;5;13m"
	ClrDk = "\033[38;5;235m"
	ClrCd = "\033[38;5;153m"
	ClrEr = "\033[38;5;216m"
	ClrBt = "\033[38;5;229m"
	Reset = "\033[0m"
	// LOG PREFIX
	Success = "success"
	Failed  = "failed"
	Inform  = "info"
	Null    = "NULL"
	FATAL   = "fatal"
)

var success = ClrOk + Success + ClrDk
var dm = Reset + "➜" + Start
var failed = ClrWn + Failed + ClrDk
var fatal = ClrWn + FATAL + ClrDk
var null = ClrWn + Null + ClrDk
var info = ClrCd + Inform + ClrDk

func Err_(msg string, sub string, err error) {
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
	log.Printf("%s %s %s %s %s\n", success, r, f, dm, p)
	return
}

func Err(r string, f string, err error) {
	log.Printf("%s %s %s %s %v\n", failed, r, f, dm, err)
	return
}

func ErrLog(r string, f string, err error) {
	if err != nil {
		log.Printf("%s %s %s %s %v\n", failed, r, f, dm, err)
		return
	}
}

func Fatal(r string, f string, err error) {
	if err != nil {
		log.Printf("%s %s %s %s %v\n", fatal, r, f, dm, err)
	}
	return
}

func OkLog(r string, f string, p interface{}, err error) {
	if err == nil {
		log.Printf("%s %s %s %s %s %v\n", success, r, f, dm, p, err)
	}
	return
}

func Info(r string, f string, p interface{}) {
	log.Printf("%s %s %s %s %s\n", info, r, f, dm, p)
}

func NullLog(r string, f string, err error) {
	log.Printf("%s %s %s %s %v\n", null, r, f, dm, err)
}

func JsonResponse(w http.ResponseWriter, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	ErrHandler(w, err)

	_, err = w.Write(jsonData)
	ErrHandler(w, err)
}

func s() string {
	return strings.ToLower(fmt.Sprintf("%04x", rand.Intn(0x10000))[1:])
}

func Guid() string {
	return fmt.Sprintf("%s%s-%s-%s-%s-%s%s%s",
		s(), s(), s(), s(), s(), s(), s(), s())
}

func RandIdx(n int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rand.New(src)
	return rand.Intn(n)
}
