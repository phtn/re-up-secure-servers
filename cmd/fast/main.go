package main

import (
	"fast/api"
	"fast/config"
	"fast/internal/models"
	"fast/internal/rdb"
	"fast/internal/repository"
	"fast/pkg/utils"
	"fmt"
	"net/http"
)

func main() {

	addr := config.LoadConfig().Addr

	mux := http.NewServeMux()

	middlewares := []api.Middleware{
		api.AuthMiddleware,
		api.CorsMiddleware,
	}

	mux.HandleFunc(api.AuthRootPath, api.Chain(api.Rdbc, middlewares...))
	mux.HandleFunc(api.GetUserPath, api.Chain(api.GetUser, middlewares...))
	mux.HandleFunc(api.CreateTokenPath, api.Chain(api.CreateToken, middlewares...))
	mux.HandleFunc(api.VerifyIdTokenPath, api.Chain(api.VerifyIdToken, middlewares...))
	mux.HandleFunc(api.VerifyAuthKeyPath, api.Chain(api.VerifyAuthKey, middlewares...))

	// DEV-ROUTES
	mux.HandleFunc(api.DevSetPath, api.Chain(api.DevSet, middlewares...))
	mux.HandleFunc(api.DevGetPath, api.Chain(api.DevGet, middlewares...))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println(repository.Code + "\n     ⟢   ╭" + repository.Dark + " ╮" + repository.Reset + "     𝗿𝗲-𝘂𝗽.𝗽𝗵 " + repository.Code)
	fmt.Println(repository.Reset + " ⟢     ╭◜" + repository.Code + "╰" + repository.Black + "⟜" + repository.Dark + "╯" + repository.Dark + "◝╮" + repository.Code + "   𝚜𝚎𝚌𝚞𝚛𝚎 ⛌ 𝚜𝚎𝚛𝚟𝚎𝚛𝚜" + repository.Start)
	fmt.Println("")

	// TEST //
	models.Ping()
	rdb.Ping()
	// END TEST //

	// SERVER START
	utils.Fatal("server", "boot", server.ListenAndServe())
	utils.Ok("server", "boot", "system-online")
}

/*
   ╭ ╮
 ╭◜╰ ╯◝╮
𝚜𝚎𝚌𝚞𝚛𝚎
𝚜𝚎𝚛𝚟𝚎𝚛𝚜
𝗿𝗲-𝘂𝗽.𝗽𝗵
𝗐 𝖾 𝖻  𝗌 𝖾 𝗋 𝗏 𝗂 𝖼 𝖾 𝗌
◜頻ꔷ⏠⛃◝◟◞─⛌྾

// fmt.Println(repository.Code + "\n ᜍ      ᜁ      ᜂ     ᜉ      ᜑ")
*/
