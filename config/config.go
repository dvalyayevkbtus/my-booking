package config

import (
	"encoding/json"
	"os"
)

var ENV_VAR = "BOOKING_CONFIG_PATH"

type Config struct {
	DB      DBConfig      `json:"db"`
	Payment PaymentConfig `json:"payment"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PaymentConfig struct {
	URL string `json:"url"`
}

func InitConfig() (Config, error) {
	PATH := os.Getenv(ENV_VAR)
	dat, err := os.ReadFile(PATH)
	if err != nil {
		return Config{}, err
	}
	var result Config
	err = json.Unmarshal(dat, &result)
	return result, err
}
