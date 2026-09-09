package utils

import (
	"frp-auth/model"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig() model.Config {
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	var cfg model.Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
