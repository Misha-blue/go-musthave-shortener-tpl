package settings

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL      string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath  string `env:"FILE_STORAGE_PATH" envDefault:"storage"`
}

func SetupConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	readFlags(&cfg)
	return &cfg, err
}

func readFlags(cfg *Config) {
	flag.StringVar(&cfg.ServerAdress, "a", cfg.ServerAdress, "port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.StoragePath, "f", cfg.StoragePath, "file storage")
	flag.Parse()
}
