package main

import (
	"fast/api"
	"fast/config"
	"fast/internal/repository"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(repository.Code + "\n ᜍ      ᜁ      ᜂ     ᜉ      ᜑ")

	addr := config.LoadConfig().ServerAddress

	http.HandleFunc("/auth", api.XXX)
	http.HandleFunc("/auth/getUser", api.GetUser)
	http.HandleFunc("/auth/createToken", api.CreateToken)
	http.HandleFunc("/auth/verifyIdToken", api.VerifyIDToken)

	// SERVER START
	// fmt.Println(repository.ColorCode + "𝙍𝙀-𝙐𝙋.𝙋𝙃 𝖢𝖫𝖮𝖴𝖣 𝖲𝖤𝖱𝖵𝖨𝖢𝖤𝖲" + repository.ColorLogStart)
	fmt.Println(repository.Reset + " 𝙧𝙚-𝙪𝙥.𝙥𝙝 " + repository.Code + " secure" + repository.Dark + "྾" + repository.Code + "servers" + repository.Dark + "  (v1)" + repository.Start)
	fmt.Println("")
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		fmt.Println(repository.Error + "server failed to start")
	}
}

/*
྾
ᜑ
𝙧𝙚-𝙪𝙥.𝙥𝙝  𝗐 𝖾 𝖻  𝗌 𝖾 𝗋 𝗏 𝗂 𝖼 𝖾 𝗌
ꔷ
頻
*/
