package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"

	"github.com/go-chi/chi"
)

type Handler struct {
	repositorier *repository.Repositorier
}

type URLPostRequest struct {
	URL string `json:"url"`
}

type URLResponseRequest struct {
	Result string `json:"result"`
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

func (handler *Handler) HandleURLJsonPostRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request := URLPostRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := handler.repositorier.Store(request.URL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := URLResponseRequest{"http://localhost:8080/" + shortURL}

	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(response)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
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
