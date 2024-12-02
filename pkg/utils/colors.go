package utils

import "github.com/muesli/termenv"

var (
	clr     = termenv.ColorProfile()
	ash     = clr.Color("#B8B4AC") //"#E1CI8E" //"#C0B6A6"
	glow    = clr.Color("#E1CI8E") //"#C0B6A6"
	raptor  = clr.Color("#F1C682") //"#C0B6A6"
	gray    = clr.Color("#1e293b") // 800
	light   = clr.Color("#4b5563") // 600
	grey    = clr.Color("#9ca3af") // 400
	sky     = clr.Color("#7dd3fc") // 500
	ice     = clr.Color("#a5f3fc")
	slate   = clr.Color("#64748b") // 500
	dsky    = clr.Color("#0369a1") // 700
	warn    = clr.Color("#fdba74") // 300
	null    = clr.Color("#ffedd5") // 200
	fail    = clr.Color("#f472b6") // 400
	good    = clr.Color("#10b981") // 500
	info    = clr.Color("#38bdf8") // 400
	debug   = clr.Color("#6366f1") // 500
	fuchsia = clr.Color("#f0abfc") // 300
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
func Ash(s string, t int) string {
	if t == 0 {
		return Fg(s, ash)
	}
	return Bg(s, ash)
}
func Glow(s string, t int) string {
	if t == 0 {
		return Fg(s, glow)
	}
	return Bg(s, glow)
}
func Raptor(s string, t int) string {
	if t == 0 {
		return Fg(s, raptor)
	}
	return Bg(s, raptor)
}
func Dev(s string, t int) string {
	if t == 0 {
		return Fg(s, debug)
	}
	return Bg(s, debug)
}
func Fuchsia(s string, t int) string {
	if t == 0 {
		return Fg(s, fuchsia)
	}
	return Bg(s, fuchsia)
}
