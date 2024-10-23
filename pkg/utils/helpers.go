package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muesli/termenv"
)

const (
	// COLORS
	Start     = "\033[38;5;235m"
	ClrOk     = "\033[38;5;35m"
	ClrWn     = "\033[38;5;216m"
	ClrNl     = "\033[38;5;175m"
	ClrDk     = "\033[38;5;248m"
	ClrCd     = "\033[38;5;153m"
	ClrEr     = "\033[38;5;168m"
	ClrBt     = "\033[38;5;59m"
	Reset     = "\033[250m"
	indicator = "⬤"
)

var success = ClrGood(indicator, 0)
var failed = ClrFail(indicator, 0)
var fatal = ClrFail(indicator, 0)
var warning = ClrWarn(indicator, 0)
var nullv = ClrWarn(indicator, 0)
var inform = ClrInfo(indicator, 0)
var div = Gray("ꔷꔷ", 0)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type Console struct {
	output  *termenv.Output
	console map[string]*log.Logger
}
type LogArgs struct {
	r string
	f string
	v interface{}
}

var (
	L = NewConsole()
	// FIBER STATUS
	Accepted     = fiber.StatusAccepted
	OK           = fiber.StatusOK
	Unauthorized = fiber.StatusUnauthorized
	BadRequest   = fiber.StatusBadRequest
)

type JsonData struct {
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Response(data interface{}, err error, message string) fiber.Map {
	response := JsonData{
		Data:    data,
		Error:   err,
		Message: message,
	}
	return fiber.Map{
		"data": response,
	}
}

var (
	format = "%s %s %s %v"
)

func NewArgs(r string, f string, v interface{}) *LogArgs {
	return &LogArgs{r: r, f: f, v: v}
}

func NewConsole() *Console {
	output := termenv.NewOutput(os.Stdout)

	console := make(map[string]*log.Logger)
	console["info"] = log.New(os.Stdout, "INFO: ", 0)
	console["good"] = log.New(os.Stdout, "GOOD: ", 0)
	console["warn"] = log.New(os.Stdout, "WARN: ", 0)
	console["null"] = log.New(os.Stdout, "NULL: ", 0)
	console["fail"] = log.New(os.Stdout, "FAIL: ", 0)
	console["fatal"] = log.New(os.Stdout, "FATAL: ", 0)

	return &Console{output, console}
}

func (l *Console) colorize(s string, c termenv.Color) string {
	return l.output.String(s).Foreground(c).String()
}

func (l *Console) formatTime() string {
	now := time.Now()
	return l.colorize(now.Format("2006-01-02 03:04:05"), light)
}

func (l *Console) formatLevel(level string) string {
	switch level {
	case "info":
		return l.colorize(indicator, info)
	case "debug":
		return l.colorize(indicator, debug)
	case "good":
		return l.colorize(indicator, good)
	case "fail":
		return l.colorize(indicator, fail)
	case "warn":
		return l.colorize(indicator, warn)
	case "null":
		return l.colorize(indicator, null)
	case "fatal":
		return l.colorize(indicator, fail)
	default:
		return level
	}
}

func (l *Console) log(level string, r string, f string, a ...interface{}) {
	msg := fmt.Sprintf(format, r, f, div, a)

	logEntry := fmt.Sprintf("%s %s %s\n", l.formatTime(), l.formatLevel(level), Grey(msg, 0))

	if level == "fail" {
		fmt.Fprintf(os.Stderr, logEntry)
	} else {
		fmt.Fprintf(os.Stdout, logEntry)
	}
}

func (l *Console) Info(r string, f string, a ...interface{}) {
	l.log("info", r, f, a...)
}

func (l *Console) Debug(r string, f string, a ...interface{}) {
	l.log("debug", r, f, a...)
}

func (l *Console) Good(r string, f string, a ...interface{}) {
	if a[len(a)-1] == nil {
		l.log("good", r, f, a...)
		return
	}
}

func (l *Console) Fail(r string, f string, a ...interface{}) {
	if a[len(a)-1] != nil {
		l.log("fail", r, f, a...)
		return
	}
}

func (l *Console) Warn(r string, f string, a ...interface{}) {
	if a[len(a)-1] != nil {
		l.log("warn", r, f, a...)
		return
	}
}

func (l *Console) Null(r string, f string, a ...interface{}) {
	l.log("null", r, f, a...)
}

func (l *Console) Fatal(r string, f string, a ...interface{}) {
	l.log("fatal", r, f, a...)
}

var (
	prefix = Gray(log.Prefix(), 0)
)

func FiberResponse(c *fiber.Ctx, status int, err error, data JsonData) error {
	if err != nil {
		L.Fail(strconv.Itoa(status), data.Message, data.Data, err)
		return c.Status(status).JSON(Response(data, err, data.Message))
	}
	L.Good(strconv.Itoa(status), data.Message, data.Data)
	return c.Status(status).JSON(data)
}

// ➜
func ErrHandler(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return err
}

func Err(r string, f string, err error) {
	log.Printf("%s %s %s %s %v\n", failed, r, f, div, err)
	return
}

func ErrLog(r string, f string, err error) {
	if err != nil {
		log.Printf("%s %s %s %s %v\n", failed, r, f, div, err)
		return
	}
}

func NoRowsErrLog(r string, f string, err error) error {
	if err != nil {
		log.Printf("%s %s %s %s %v\n", failed, r, f, div, err)
		if err == sql.ErrNoRows {
			log.Printf("%s %s %s %s %v\n", failed, r, f, div, err)
			return nil
		}
		return err
	}
	return err
}

func HttpError(w http.ResponseWriter, m string, err error, code int) {
	if err != nil {
		v := ErrorResponse{
			Status:  "error",
			Message: m,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(v)
		http.Error(w, m, code)
		Err("http", m, err)
		return
	}
}

func HttpErr(w http.ResponseWriter, m string, err error, code int) {
	v := ErrorResponse{
		Status:  "error",
		Message: m,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
	http.Error(w, m, code)
	Err("http", m, err)
	return
}

func Warn(r string, f string, p interface{}) {
	log.Printf("%s %s %s %s %s\n", warning, r, f, div, p)
	return
}

func WarnLog(r string, f string, p interface{}, err error) {
	format := prefix + "%s %s %s %s %s %v\n"
	if err == nil {
		log.Printf(format, warning, r, f, div, p, err)
		return
	}
}

func Ok(r string, f string, p interface{}) {
	log.Printf("%s %s %s %s %s\n", success, Grey(r, 0), f, div, p)
	return
}

func OkLog(r string, f string, p interface{}, err error) {
	if err == nil {
		log.Printf("%s %s %s %s %s %v\n", success, Light(r, 0), f, div, p, err)
		return
	}
}

func Info(r string, f string, p interface{}) {
	log.Printf("%s %s %s %s %s\n", inform, r, f, div, p)
	return
}

func PostMethodOnly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Info("method", "verifyIdToken", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetMethodOnly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Info("method", "verifyIdToken", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Fatal(r string, f string, err error) {
	if err != nil {
		log.Printf("%s %s %s %s %v\n", fatal, r, f, div, err)
		return
	}
}

func NullLog(r string, f string, err error) {
	log.Printf("%s %s %s %s %v\n", null, r, f, div, err)
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
