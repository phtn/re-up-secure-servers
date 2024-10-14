package main

import (
	"fast/api"
	"fast/config"
	"fast/internal/repository"
	"fast/pkg/utils"
	"fmt"
	"net/http"
)

func main() {

	addr := config.LoadConfig().Addr

	http.HandleFunc(api.AuthRootPath, api.RDBC)
	http.HandleFunc(api.GetUserPath, api.GetUser)
	http.HandleFunc(api.CreateTokenPath, api.CreateToken)
	http.HandleFunc(api.VerifyIdTokenPath, api.VerifyIdToken)
	http.HandleFunc(api.VerifyAuthKeyPath, api.VerifyAuthKey)

	// DEV-ROUTES
	http.HandleFunc(api.DevSetPath, api.DevSet)
	http.HandleFunc(api.DevGetPath, api.DevGet)

	utils.Ok("server", "boot", "system-online")
	fmt.Println("")
	// SERVER START
	fmt.Println(repository.Code + "\n     ⟢   ╭" + repository.Dark + " ╮" + repository.Reset + "     𝗿𝗲-𝘂𝗽.𝗽𝗵 " + repository.Code)
	fmt.Println(repository.Reset + " ⟢     ╭◜" + repository.Code + "╰" + repository.Black + "⟜" + repository.Dark + "╯" + repository.Dark + "◝╮" + repository.Code + "   𝚜𝚎𝚌𝚞𝚛𝚎 ⛌ 𝚜𝚎𝚛𝚟𝚎𝚛𝚜" + repository.Start)
	fmt.Println("")
	err := http.ListenAndServe(addr, nil)

	utils.Fatal("server", "boot", err)
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
