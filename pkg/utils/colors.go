package utils

import "github.com/muesli/termenv"

var (
	clr   = termenv.ColorProfile()
	gray  = clr.Color("#1e293b") // 800
	light = clr.Color("#4b5563") // 600
	grey  = clr.Color("#9ca3af") // 400
	sky   = clr.Color("#7dd3fc") // 500
	ice   = clr.Color("#a5f3fc")
	slate = clr.Color("#64748b") // 500
	dsky  = clr.Color("#0369a1") // 700
	warn  = clr.Color("#fdba74") // 300
	null  = clr.Color("#ffedd5") // 200
	fail  = clr.Color("#f472b6") // 400
	good  = clr.Color("#10b981") // 500
	info  = clr.Color("#38bdf8") // 400
	debug = clr.Color("#6366f1") // 500
)

func Colorize(s string, c termenv.Color) string {
	cstr := clr.String(s)
	return cstr.Foreground(c).String()
}

func Fg(s string, c termenv.Color) string {
	cstr := clr.String(s)
	return cstr.Foreground(c).String()
}
func Bg(s string, c termenv.Color) string {
	cstr := clr.String(s)
	return cstr.Background(c).String()
}

func Gray(s string, t int) string {
	if t == 0 {
		return Fg(s, gray)
	}
	return Bg(s, gray)
}
func Sky(s string, t int) string {
	if t == 0 {
		return Fg(s, sky)
	}
	return Bg(s, sky)
}

func DSky(s string, t int) string {
	if t == 0 {
		return Fg(s, dsky)
	}
	return Bg(s, dsky)
}
func Ice(s string, t int) string {
	if t == 0 {
		return Fg(s, ice)
	}
	return Bg(s, ice)
}
func Light(s string, t int) string {
	if t == 0 {
		return Fg(s, light)
	}
	return Bg(s, light)
}

func Slate(s string, t int) string {
	if t == 0 {
		return Fg(s, slate)
	}
	return Bg(s, slate)
}

func ClrWarn(s string, t int) string {
	if t == 0 {
		return Fg(s, warn)
	}
	return Bg(s, warn)
}
func ClrGood(s string, t int) string {
	if t == 0 {
		return Fg(s, good)
	}
	return Bg(s, good)
}
func ClrFail(s string, t int) string {
	if t == 0 {
		return Fg(s, fail)
	}
	return Bg(s, fail)
}
func ClrInfo(s string, t int) string {
	if t == 0 {
		return Fg(s, info)
	}
	return Bg(s, info)
}
func Grey(s string, t int) string {
	if t == 0 {
		return Fg(s, grey)
	}
	return Bg(s, grey)
}
