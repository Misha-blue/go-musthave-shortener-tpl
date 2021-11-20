package main

import (
	"log"
	"net/http"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HandleURLRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
