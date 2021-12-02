package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Misha-blue/go-musthave-shortener-tpl/cmd/shortener/server"
	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAdress string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
}

func main() {
	cfg := SetupConfig()

	repository := repository.New()
	handler := handlers.New(&repository, cfg.BaseURL)
	server := server.New(handler, cfg.ServerAdress)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		oscall := <-done
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := server.Run(ctx); err != nil {
		log.Printf("failed to run server:+%v\n", err)
	}
}

func SetupConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Print(err)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}

	if cfg.ServerAdress == "" {
		cfg.ServerAdress = ":8080"
	}

	return cfg
}
