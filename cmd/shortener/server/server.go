package server

import (
	"context"
	"log"
	"net/http"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi"
)

type Server struct{}

func (s *Server) Run(ctx context.Context) (*http.Server, error) {
	router := chi.NewRouter()
	handler := handlers.New(repository.New())

	router.Get("/{shortURL}", handler.HandleURLGetRequest)
	router.Post("/", handler.HandleURLPostRequest)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server didn't run: %s\n", err)
		}
	}()

	return &server, nil
}

func (s *Server) Close(server *http.Server, ctx context.Context) {
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	log.Print("Server exited properly")
}
