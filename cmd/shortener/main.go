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
	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository/file"
)

func main() {
	cfg, err := settings.SetupConfig()
	if err != nil {
		log.Printf("Failed to read environment settings:+%v\n", err)
	}

	log.Print(cfg)
	storage, err := file.New(cfg.StoragePath + "/fileStorage.txt")

	if err != nil {
		log.Printf("Failed to create storage:+%v\n", err)
	}
	defer storage.Close()

	repository := repository.New(storage)
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
