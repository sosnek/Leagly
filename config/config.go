package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	Token  string
	ApiKey string

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token  string `json : "Token"`
	ApiKey string `json : "ApiKey"`
}

func ReadConfig() error {

	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("config/config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	ApiKey = config.ApiKey

	if Token == "" {
		Token = os.Getenv("DiscordKey")
	}
	if ApiKey == "" {
		ApiKey = os.Getenv("RiotKey")
	}

	return nil
}
