package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi"
)

type Server struct{}

func (s *Server) Run(ctx context.Context) error {
	router := chi.NewRouter()
	handler := handlers.New(repository.New())

	router.Get("/{shortURL}", handler.HandleURLGetRequest)
	router.Post("/", handler.HandleURLPostRequest)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	log.Print("Server started")

	<-done
	log.Print("Server stopped")

	return server.Shutdown(ctx)
}
