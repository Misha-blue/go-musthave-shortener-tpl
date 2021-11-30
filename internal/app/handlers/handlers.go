package handlers

import (
	"io"
	"net/http"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"

	"github.com/go-chi/chi"
)

type Handler struct {
	repositorier *repository.Repositorier
}

func New(repositorier *repository.Repositorier) *Handler {
	return &Handler{repositorier: repositorier}
}

func (handler *Handler) HandleURLPostRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := handler.repositorier.Store(string(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shortURL))
}

func (handler *Handler) HandleURLGetRequest(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	url, err := handler.repositorier.Load(shortURL)

	if err == nil {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
