package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi"
)

type Server struct{}

func (s *Server) Run(ctx context.Context) error {
	handler := chi.NewRouter()

	handler.Get("/{shortURL}", handlers.HandleURLGetRequest)
	handler.Post("/", handlers.HandleURLPostRequest)

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server.ListenAndServe()
	}()

	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	return server.Shutdown(ctx)
}
