package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Misha-blue/go-musthave-shortener-tpl/cmd/shortener/server"
	"github.com/Misha-blue/go-musthave-shortener-tpl/cmd/shortener/settings"
	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"
)

func main() {
	cfg, err := settings.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to read environment settings:+%v\n", err)
	}

	repository, err := repository.New(cfg.StoragePath)
	if err != nil {
		log.Fatalf("Failed to create storage:+%v\n", err)
	}

	handler := handlers.New(repository, cfg.BaseURL)
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
