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
	Start = "\033[38;5;235m"
	ClrOk = "\033[38;5;35m"
	ClrWn = "\033[38;5;216m"
	ClrNl = "\033[38;5;175m"
	ClrDk = "\033[38;5;248m"
	ClrCd = "\033[38;5;153m"
	ClrEr = "\033[38;5;168m"
	ClrBt = "\033[38;5;59m"
	Reset = "\033[250m"
	// LOG PREFIX
	Success = "success"
	Warning = "warning"
	Failed  = "failed "
	Inform  = "info   "
	Null    = "NULL   "
	FATAL   = "fatal  "
)

var success = ClrOk + Success + ClrDk
var failed = ClrEr + Failed + ClrDk
var fatal = ClrEr + FATAL + ClrDk
var warning = ClrWn + Warning + ClrDk
var null = ClrNl + Null + ClrDk
var info = ClrCd + Inform + ClrDk
var dm = ClrBt + "ꔷ" + Start

// ➜
func ErrHandler(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return err
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

func HttpError(w http.ResponseWriter, m string, err error) {
	if err != nil {
		http.Error(w, m, http.StatusInternalServerError)
	}
}

func Warn(r string, f string, p interface{}) {
	log.Printf("%s %s %s %s %s\n", warning, r, f, dm, p)
	return
}

func Ok(r string, f string, p interface{}) {
	log.Printf("%s %s %s %s %s\n", success, r, f, dm, p)
	return
}

func PostMethodOnly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Info("method", "verifyIdToken", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
