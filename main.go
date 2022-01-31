package main

import (
	"Leagly/bot"
	"Leagly/config"
	"fmt"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	bot.ConnectToDiscord()
}
