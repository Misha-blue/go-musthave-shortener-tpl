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
)

func main() {
	repository := repository.New()
	handler := handlers.New(repository)
	server := server.New(handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	oscall := <-done
	log.Printf("system call:%+v", oscall)

	if err := server.Run(ctx); err != nil {
		log.Printf("failed to run server:+%v\n", err)
	}
}
