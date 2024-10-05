package main

import (
	"fast/api"
	"fast/config"
	"fast/internal/repository"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(repository.ColorCode + "\n ᜍ      ᜁ      ᜂ     ᜉ      ᜑ")

	addr := config.LoadConfig().ServerAddress

	http.HandleFunc("/auth", api.XXX)
	http.HandleFunc("/auth/createToken", api.CreateToken)
	http.HandleFunc("/auth/getUser", api.GetUser)

	// SERVER START
	// fmt.Println(repository.ColorCode + "𝙍𝙀-𝙐𝙋.𝙋𝙃 𝖢𝖫𝖮𝖴𝖣 𝖲𝖤𝖱𝖵𝖨𝖢𝖤𝖲" + repository.ColorLogStart)
	fmt.Println(repository.ColorReset + " 𝙧𝙚-𝙪𝙥.𝙥𝙝 " + repository.ColorResp + " secure" + repository.ColorDark + "྾" + repository.ColorResp + "servers" + repository.ColorDark + "  (v1)" + repository.ColorLogStart)
	fmt.Println("")
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		fmt.Println("failed to start server", err)
	}
}

/*

྾
ᜑ
ꔷ
頻
頻
𝙧𝙚-𝙪𝙥.𝙥𝙝  𝗐 𝖾 𝖻  𝗌 𝖾 𝗋 𝗏 𝗂 𝖼 𝖾 𝗌
𝖶
𝖡
𝙐𝙝𝙃𝙋









*/
