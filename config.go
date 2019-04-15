package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ServiceUrl                   string
	MysqlHost                    string
	MysqlDb                      string
	MysqlUser                    string
	MysqlPassword                string
	MySqlResidentalCompoundTable string
	MySqlCorpusTable             string
	MySqlApartmentTable          string
	MySqlApartmentView           string
}

func LoadConfig() *Config {
	file, _ := os.Open("config.json")

	decoder := json.NewDecoder(file)
	config := new(Config)
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal("invalid config file")
	}
	return config
}
