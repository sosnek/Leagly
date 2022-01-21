package config

import (
	"encoding/json"
	"fmt"       //used to print errors majorly.
	"io/ioutil" //it will be used to help us read our config.json file.
)

var (
	Token     string //To store value of Token from config.json .
	BotPrefix string // To store value of BotPrefix from config.json.
	Region    string
	ApiKey    string

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token     string `json : "Token"`
	BotPrefix string `json : "BotPrefix"`
	Region    string `json : "Region"`
	ApiKey    string `json : "ApiKey"`
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
	BotPrefix = config.BotPrefix
	Region = config.Region
	ApiKey = config.ApiKey

	return nil
}
