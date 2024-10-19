package repository

type Color string

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

type Method string

const (
	Get    Method = "get"
	Post   Method = "post"
	Put    Method = "put"
	Delete Method = "delete"
	Patch  Method = "patch"
)
