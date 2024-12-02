package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/muesli/termenv"
)

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
	L      = NewConsole()
	prefix = Gray(log.Prefix(), 0)
)

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
	return l.colorize(now.Format("03:04"), light)
}

func (l *Console) formatLevel(level string) string {
	switch level {
	case "info":
		return l.colorize(indicator, info)
	case "debug":
		return l.colorize(d_indicator, sky)
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

	input := []string{l.formatTime(), l.formatLevel(level), Ash(msg, 0)}
	logEntry := fmt.Sprintf("%s\n", OpenRect(input))

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
	l.log("info", r, f, "OK")
	return
}

func (l *Console) FailR(r string, f string, a ...interface{}) (interface{}, error) {
	if a[len(a)-1] != nil {
		l.log("fail", r, f, a...)
		return nil, errors.New(a[len(a)-1].(string))
	}
	l.log("good", r, f, "OK")
	return a[1:], nil
}

func FailRII[E any](r string, f string, a []E) ([]E, error) {
	if a != nil {
		e := fmt.Errorf("%v", a)
		L.log("fail", r, f, e)
		return nil, e
	}
	L.log("good", r, f, "OK")
	return a, nil
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
