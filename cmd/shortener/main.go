package main

import (
	"context"
	"log"

	"github.com/Misha-blue/go-musthave-shortener-tpl/cmd/shortener/server"
)

func main() {
	server := &server.Server{}
	log.Fatal(server.Run(context.Background()))
}
