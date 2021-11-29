package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Misha-blue/go-musthave-shortener-tpl/cmd/shortener/server"
)

func main() {
	server := &server.Server{}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv, _ := server.Run(context.Background())
	log.Print("Server started")

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Close(srv, ctx)
}
