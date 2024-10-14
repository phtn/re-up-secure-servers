package repository

type Color string

const (
	Start   = "\033[38;5;60m"
	Success = "\033[38;5;150m"
	Warn    = "\033[38;5;13m"
	Dark    = "\033[38;5;235m"
	Black   = "\033[38;5;233m"
	Code    = "\033[38;5;153m"
	Error   = "\033[38;5;216m"
	Bright  = "\033[38;5;229m"
	Reset   = "\033[0m"
)

type Method string

const (
	Get    Method = "get"
	Post   Method = "post"
	Put    Method = "put"
	Delete Method = "delete"
	Patch  Method = "patch"
)
