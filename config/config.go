package config

import (
	"encoding/json"
	"fmt"       //used to print errors majorly.
	"io/ioutil" //it will be used to help us read our config.json file.
	"os"
)

var (
	Token         string //To store value of Token from config.json .
	EncryptionKey string // To store value of BotPrefix from config.json.
	ApiKey        string

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token         string `json : "Token"`
	EncryptionKey string `json : "encryptionKey"`
	ApiKey        string `json : "ApiKey"`
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
	EncryptionKey = config.EncryptionKey
	ApiKey = config.ApiKey

	if Token == "" {
		Token = os.Getenv("DiscordKey")
	}
	if ApiKey == "" {
		ApiKey = os.Getenv("RiotKey")
	}
	if EncryptionKey == "" {
		ApiKey = os.Getenv("EncryptionKey")
	}

	return nil
}
