package main

import (
	"fast/api"
	"fast/config"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/pkg/utils"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/muesli/termenv"
)

var (
	tclr = termenv.ColorProfile()
)

func main() {

	addr := config.LoadConfig().Addr

	mux := http.NewServeMux()

	middlewares := []api.Middleware{
		api.AuthMiddleware,
		api.CorsMiddleware,
	}
	admin_middlewares := []api.Middleware{
		api.AuthMiddleware,
		api.CorsMiddleware,
		api.AdminClaimsMiddleware,
	}
	withClaims := append(middlewares, api.ClaimsMiddleware)
	withAdmin := append(middlewares, api.AdminClaimsMiddleware)

	mux.HandleFunc(api.AuthPath, api.Chain(api.DbCheck, middlewares...))
	mux.HandleFunc(api.GetUserPath, api.Chain(api.GetUser, middlewares...))
	mux.HandleFunc(api.CreateTokenPath, api.Chain(api.CreateToken, middlewares...))
	mux.HandleFunc(api.VerifyIdTokenPath, api.Chain(api.VerifyIdToken, middlewares...))
	mux.HandleFunc(api.VerifyAuthKeyPath, api.Chain(api.VerifyAuthKey, middlewares...))

	// WITH CLAIMS
	mux.HandleFunc(api.CustomClaimsPath, api.Chain(api.CreateCustomClaims, withClaims...))
	// WITH ADMIN
	mux.HandleFunc(api.AdminPath, api.Chain(api.CheckAdminAuthority, admin_middlewares...))
	mux.HandleFunc(api.AdminClaimsPath, api.Chain(api.CreateAdminClaims, withAdmin...))

	// DEV-ROUTES
	mux.HandleFunc(api.DevSetPath, api.Chain(api.DevSet, middlewares...))
	mux.HandleFunc(api.DevGetPath, api.Chain(api.DevGet, middlewares...))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// splash()
	MkOne()

	// TEST //
	models.PsqlHealth()
	rdb.RedisHealth()
	// END TEST //

	// SERVER START
	err := server.ListenAndServe()
	utils.Fatal("serve", "boot", err)
	utils.OkLog("serve", "boot", "system-online", err)
}

var (
	gr8 = tclr.Color("#374151")
)

func drawBorder(l int, t int) {

	rt := []string{"╭", "╮"}
	rb := []string{"╰", "╯"}

	// createLine( )
	c := rt
	if t == 1 {
		c = rb
	}

	fmt.Printf(Colorize(c[0], gr8))
	progress := "──"
	for range l {
		time.Sleep(25 * time.Millisecond)
		fmt.Printf(Colorize(progress, gr8))
	}
	fmt.Println(Colorize(c[1], gr8))
}

func renderContent(c string, l int) {

	p := l - countVis(c)
	if p < 0 {
		p = p * -1
	}

	v := Colorize("│", gr8)
	ws := strings.Repeat(" ", l-countVis(c))

	// fmt.Printf(" %v ", (l-countVis(c))/2)
	fmt.Println(v + c + ws + v)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func Colorize(s string, c termenv.Color) string {
	cstr := tclr.String(s)
	return cstr.Foreground(c).String()
}
func splash() {

	clearScreen()

	rose := tclr.Color("#fb7185")

	// width, _, err := term.GetSize(int(os.Stdout.Fd()))
	// utils.ErrLog("line", "get-size", err)

	// contents := []string{firstrow, secndrow, thirdrow}

	// n := (width / 6) - 2

	// for i := range len(contents) {
	// 	renderContent(contents[i], n/2)
	// }
	// createLine(n, "╰", "╯")

	// fmt.Println(strings.Count(firstrow, ""), len(secndrow), strings.Count(thirdrow, ""))
	fmt.Println(Colorize("", rose))

}

func countVis(input string) int {
	// Regular expression to match ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

	// Remove all ANSI escape sequences from the input string
	cleanString := ansiRegex.ReplaceAllString(input, "")

	// Count the number of characters in the cleaned string
	return utf8.RuneCountInString(cleanString)
}

func MkOne() {

	clearScreen()

	f := "    ⟢     ╭ ╮"
	s := " ⟢      ╭◜╰ ╯◝╮"
	t := "   ⟢       ◌"

	drawBorder(9, 0)
	for c := range map[string]interface{}{f: f, s: s, t: t} {
		renderContent(c, 18)
	}
	drawBorder(9, 1)
}
