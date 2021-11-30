package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Misha-blue/go-musthave-shortener-tpl/cmd/shortener/server"
)

func main() {
	server := &server.Server{}
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
