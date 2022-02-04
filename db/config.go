package db

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Password string `json:"password"`
		Username string `json:"username"`
		Db       string `json:"db"`
	} `json:"database"`
	Telegram struct {
		Token string
		Chat  string
	}
	Sleep int  `json:"sleep"`
	Debug bool `json:"debug"`
}

func LoadConfiguration(file string) Config {
	var config Config

	pwd, _ := os.Getwd()
	configFile, err := os.Open(pwd + "/" + file)
	defer configFile.Close()
	if err != nil {
		panic(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
