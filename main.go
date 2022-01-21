package main

import (
	"Leagly/bot"
	"Leagly/config"
	"Leagly/query"
	"fmt"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	query.InitializedChampStruct()
	bot.ConnectToDiscord()
}
