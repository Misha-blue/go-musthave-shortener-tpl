package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	handler       *handlers.Handler
	serverAddress string
}

func New(handler *handlers.Handler, serverAddress string) *Server {
	return &Server{
		handler:       handler,
		serverAddress: serverAddress,
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	serverCtx, cancel := context.WithCancel(ctx)

	router := chi.NewRouter()

	router.Use(middleware.Compress(5))
	router.Get("/{shortURL}", s.handler.HandleURLGetRequest)
	router.Post("/", s.handler.HandleURLPostRequest)
	router.Post("/api/shorten", s.handler.HandleURLJsonPostRequest)

	server := http.Server{
		Addr:    s.serverAddress,
		Handler: router,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server didn't start: %+s\n", err)
		}
		cancel()
	}()

	log.Printf("server started")

	<-serverCtx.Done()
	log.Printf("server stopping")

	shutDownCtx, shutDownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutDownCancel()

	if err = server.Shutdown(shutDownCtx); err != nil {
		log.Printf("server shutdown failed: %+s", err)
	}

	log.Printf("server stopped properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
