package main

import (
	"fmt"
	"Leagly/bot"
	"Leagly/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
  
	bot.ConnectToDiscord()
  
}