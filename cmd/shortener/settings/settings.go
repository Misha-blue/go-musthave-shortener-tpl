package settings

import (
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL      string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func SetupConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	log.Print(&cfg)
	return &cfg, err
}
