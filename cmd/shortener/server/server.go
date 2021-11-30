package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi"
)

type Server struct{}

func (s *Server) Run(ctx context.Context) (err error) {
	router := chi.NewRouter()
	handler := handlers.New(repository.New())

	router.Get("/{shortURL}", handler.HandleURLGetRequest)
	router.Post("/", handler.HandleURLPostRequest)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server didn't run:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
