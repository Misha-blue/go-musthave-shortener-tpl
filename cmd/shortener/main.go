package main

import (
	"log"
	"net/http"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/{shortURL}", handlers.HandleURLGetRequest)
	r.Post("/", handlers.HandleURLPostRequest)

	log.Fatal(http.ListenAndServe(":8080", r))
}
